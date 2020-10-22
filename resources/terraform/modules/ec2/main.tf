locals {
  use_nat_gateway = var.force_nat_gateway || !var.use_public_ip
}

resource "aws_resourcegroups_group" "rg" {
  name = "rg-${var.name}"

  resource_query {
    query = <<JSON
{
  "ResourceTypeFilters": [
    "AWS::EC2::Subnet",
    "AWS::EC2::VPC",
    "AWS::EC2::InternetGateway",
    "AWS::EC2::NatGateway",
    "AWS::EC2::SecurityGroup",
    "AWS::EC2::Instance",
    "AWS::EC2::RouteTable"
  ],
  "TagFilters": [
    {
      "Key": "cluster_name",
      "Values": ["${var.name}"]
    }
  ]
}
JSON
  }
}

resource "aws_vpc" "awsbi_vpc" {
  cidr_block            = var.vpc_cidr_block
  instance_tenancy      = "default"
  enable_dns_support    = "true"
  enable_dns_hostnames  = "true"

  tags = {
    Name          = "vpc-${var.name}"
    cluster_name  = var.name
  }
}

resource "aws_subnet" "awsbi_private_subnet" {
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnet_private_cidr_block
  availability_zone = "${var.region}a"

  tags = {
    Name         = "subnet-private-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_subnet" "awsbi_public_subnet" {
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnet_public_cidr_block
  availability_zone = "${var.region}a"

  tags = {
    Name         = "subnet-public-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_internet_gateway" "awsbi_internet_gateway" {
  vpc_id = aws_vpc.awsbi_vpc.id

  tags = {
    Name         = "ig-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_eip" "awsbi_nat_gateway" {
  count = local.use_nat_gateway ? 1 : 0

  vpc = true
}

resource "aws_nat_gateway" "awsbi_nat_gateway" {
  count = local.use_nat_gateway ? 1 : 0

  allocation_id = aws_eip.awsbi_nat_gateway.*.id[0]
  subnet_id     = aws_subnet.awsbi_public_subnet.id

  tags = {
    Name         = "ng-${var.name}"
    cluster_name = var.name
  }

  depends_on = [ aws_internet_gateway.awsbi_internet_gateway ]
}

resource "aws_route_table" "awsbi_route_table_private" {
  count = local.use_nat_gateway ? 1 : 0

  vpc_id = aws_vpc.awsbi_vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.awsbi_nat_gateway.*.id[0]
  }

  tags = {
    Name = "rt-private-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_route_table" "awsbi_route_table_public" {
  vpc_id  = aws_vpc.awsbi_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.awsbi_internet_gateway.id
  }

  tags = {
    Name = "rt-public-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_route_table_association" "awsbi_route_association_private" {
  subnet_id      = aws_subnet.awsbi_private_subnet.id
  route_table_id = var.use_public_ip ? aws_route_table.awsbi_route_table_public.id : aws_route_table.awsbi_route_table_private.*.id[0]
}

resource "aws_route_table_association" "awsbi_route_association_public" {
  subnet_id      = aws_subnet.awsbi_public_subnet.id
  route_table_id = aws_route_table.awsbi_route_table_public.id
}

resource "aws_security_group" "awsbi_security_group" {
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
    Name = "sg-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_instance" "awsbi" {
  count = var.instance_count

  ami                         = var.ami
  instance_type               = var.instance_type
  subnet_id                   = aws_subnet.awsbi_private_subnet.id
  associate_public_ip_address = var.use_public_ip
  key_name                    = var.key_name

  root_block_device {
    volume_size = var.root_volume_size
  }

  vpc_security_group_ids = [
    aws_security_group.awsbi_security_group.id
  ]

  tags = {
    Name = var.name
    cluster_name = var.name
  }
}
