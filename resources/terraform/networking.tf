data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "awsbi_vpc" {
  cidr_block            = var.vpc_address_space
  instance_tenancy      = "default"
  enable_dns_support    = "true"
  enable_dns_hostnames  = "true"

  tags = {
    Name           = "${var.name}-vpc"
    resource_group = var.name
  }
}

resource "aws_security_group" "awsbi_security_group" {
  count   = length(var.security_groups)
  name    = var.security_groups[count.index].name
  vpc_id  = aws_vpc.awsbi_vpc.id

  dynamic "ingress" {
    for_each = var.security_groups[count.index].rules.ingress
    content {
      protocol    = ingress.value["protocol"]
      from_port   = ingress.value["from_port"]
      to_port     = ingress.value["to_port"]
      cidr_blocks = ingress.value["cidr_blocks"]
    }
  }

  dynamic "egress" {
    for_each = var.security_groups[count.index].rules.egress
    content {
      protocol    = egress.value["protocol"]
      from_port   = egress.value["from_port"]
      to_port     = egress.value["to_port"]
      cidr_blocks = egress.value["cidr_blocks"]
    }
  }

  tags = {
    Name           = "${var.name}-sg-${count.index}"
    resource_group = var.name
  }
}

# --- Public ---

resource "aws_subnet" "awsbi_public_subnet" {
  count             = length(var.subnets.public)
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnets.public[count.index].address_prefixes
  availability_zone = var.subnets.public[count.index].availability_zone == "any" ? element(data.aws_availability_zones.available.names, count.index) : var.subnets.public[count.index].availability_zone

  tags = {
    Name           = var.subnets.public[count.index].name
    resource_group = var.name
  }
}

resource "aws_internet_gateway" "awsbi_internet_gateway" {
  count  = local.use_internet_gateway ? 1 : 0
  vpc_id = aws_vpc.awsbi_vpc.id

  tags = {
    Name           = "${var.name}-ig"
    resource_group = var.name
  }
}

resource "aws_route_table" "awsbi_route_table_public" {
  count  = local.use_internet_gateway ? 1 : 0
  vpc_id = aws_vpc.awsbi_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.awsbi_internet_gateway[0].id
  }

  tags = {
    Name           = "${var.name}-rt-public"
    resource_group = var.name
  }
}

resource "aws_route_table_association" "awsbi_route_association_public" {
  count          = length(var.subnets.public)
  subnet_id      = aws_subnet.awsbi_public_subnet[count.index].id
  route_table_id = aws_route_table.awsbi_route_table_public[0].id
}

# --- Private ---


resource "aws_subnet" "awsbi_private_subnet" {
  count             = length(var.subnets.private)
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnets.private[count.index].address_prefixes
  availability_zone = var.subnets.private[count.index].availability_zone == "any" ? element(data.aws_availability_zones.available.names, count.index) : var.subnets.private[count.index].availability_zone

  tags = {
    Name           = var.subnets.private[count.index].name
    resource_group = var.name
  }
}

resource "aws_vpn_gateway" "awsbi_vpn_gw" {
  count  = local.use_virtual_private_gateway ? 1 : 0
  vpc_id = aws_vpc.awsbi_vpc.id

  tags = {
    Name = "${var.name}-vpn-gw${count.index}"
  }
}

resource "aws_eip" "awsbi_nat_gateway" {
  count = local.use_nat_gateway ? var.nat_gateway_count : 0 
  vpc   = true

  tags = {
    Name           = "${var.name}-eip${count.index}"
    resource_group = var.name
  }
}

resource "aws_nat_gateway" "awsbi_nat_gateway" {
  count         = local.use_nat_gateway ? var.nat_gateway_count : 0
  allocation_id = aws_eip.awsbi_nat_gateway[count.index].id
  subnet_id     = element(aws_subnet.awsbi_public_subnet.*.id, count.index)

  tags = {
    Name           = "${var.name}-ng${count.index}"
    resource_group = var.name
  }

  depends_on = [ aws_internet_gateway.awsbi_internet_gateway ]
}

resource "aws_route_table" "awsbi_route_table_private" {
  count  = local.use_nat_gateway ? var.nat_gateway_count : (local.use_virtual_private_gateway ? 1 : 0)
  vpc_id = aws_vpc.awsbi_vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = local.use_nat_gateway ? aws_nat_gateway.awsbi_nat_gateway[count.index].id : null
    gateway_id     = local.use_virtual_private_gateway ? aws_vpn_gateway.awsbi_vpn_gw[count.index].id : null
  }

  tags = {
    Name           = "${var.name}-rt-private${count.index}"
    resource_group = var.name
  }
}

resource "aws_route_table_association" "awsbi_route_association_private" {
  count          = local.use_nat_gateway || local.use_virtual_private_gateway ? length(var.subnets.private) : 0
  subnet_id      = aws_subnet.awsbi_private_subnet[count.index].id
  route_table_id = element(aws_route_table.awsbi_route_table_private.*.id, count.index)
}
