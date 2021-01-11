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
  name    = "${var.name}-sg"
  vpc_id  = aws_vpc.awsbi_vpc.id

  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    resource_group = var.name
  }
}

# --- Public ---

resource "aws_subnet" "awsbi_public_subnet" {
  count             = local.public_subnet_numbers
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnets.public[count.index].address_prefixes
  availability_zone = var.subnets.public[count.index].availability_zone == "any" ? element(data.aws_availability_zones.available.names, count.index) : var.subnets.public[count.index].availability_zone

  tags = {
    Name           = var.subnets.public[count.index].name
    resource_group = var.name
  }
}

resource "aws_internet_gateway" "awsbi_internet_gateway" {
  vpc_id = aws_vpc.awsbi_vpc.id

  tags = {
    Name           = "${var.name}-ig"
    resource_group = var.name
  }
}

resource "aws_route_table" "awsbi_route_table_public" {
  vpc_id  = aws_vpc.awsbi_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.awsbi_internet_gateway.id
  }

  tags = {
    Name           = "${var.name}-rt-public"
    resource_group = var.name
  }
}

resource "aws_route_table_association" "awsbi_route_association_public" {
  count          = local.public_subnet_numbers
  subnet_id      = aws_subnet.awsbi_public_subnet[count.index].id
  route_table_id = aws_route_table.awsbi_route_table_public.id
}

# --- Private ---


resource "aws_subnet" "awsbi_private_subnet" {
  count             = local.private_subnet_numbers
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnets.private[count.index].address_prefixes
  availability_zone = var.subnets.private[count.index].availability_zone == "any" ? element(data.aws_availability_zones.available.names, count.index) : var.subnets.private[count.index].availability_zone

  tags = {
    Name           = var.subnets.private[count.index].name
    resource_group = var.name
  }
}

resource "aws_eip" "awsbi_nat_gateway" {
  count = var.nat_gateway_count
  vpc   = true

  tags = {
    Name           = "${var.name}-eip${count.index}"
    resource_group = var.name
  }
}

resource "aws_nat_gateway" "awsbi_nat_gateway" {
  count         = var.nat_gateway_count
  allocation_id = aws_eip.awsbi_nat_gateway[count.index].id
  subnet_id     = element(aws_subnet.awsbi_public_subnet.*.id, count.index)

  tags = {
    Name           = "${var.name}-ng${count.index}"
    resource_group = var.name
  }

  depends_on = [ aws_internet_gateway.awsbi_internet_gateway ]
}

resource "aws_route_table" "awsbi_route_table_private" {
  count = var.nat_gateway_count
  vpc_id = aws_vpc.awsbi_vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.awsbi_nat_gateway[count.index].id
  }

  tags = {
    Name           = "${var.name}-rt-private${count.index}"
    resource_group = var.name
  }
}

resource "aws_route_table_association" "awsbi_route_association_private" {
  count          = local.private_subnet_numbers
  subnet_id      = aws_subnet.awsbi_private_subnet[count.index].id
  route_table_id = element(aws_route_table.awsbi_route_table_private.*.id, count.index)
}
