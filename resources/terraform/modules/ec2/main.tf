resource "aws_resourcegroups_group" "rg" {
  name = "rg-${var.name}"

  resource_query {
    query = <<JSON
{
  "ResourceTypeFilters": [
    "AWS::EC2::Subnet",
    "AWS::EC2::VPC",
    "AWS::EC2::InternetGateway",
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

resource "aws_subnet" "awsbi_subnet" {
  vpc_id            = aws_vpc.awsbi_vpc.id
  cidr_block        = var.subnet_cidr_block
  availability_zone = "${var.region}a"

  tags = {
    Name = "subnet-${var.name}"
    cluster_name = var.name
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

resource "aws_internet_gateway" "awsbi_internet_gateway" {
  vpc_id  = aws_vpc.awsbi_vpc.id
  tags    = {
    Name          = "internet-gateway-${var.name}"
    cluster_name  = var.name
  }
}

resource "aws_security_group" "awsbi_security_group" {
  vpc_id  = aws_vpc.awsbi_vpc.id
  tags    = {
    Name = "sg-${var.name}"
    cluster_name = var.name
  }

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
}

resource "aws_instance" "awsbi" {
  count                       = var.instance_count
  ami                         = var.ami
  instance_type               = var.instance_type
  subnet_id                   = aws_subnet.awsbi_subnet.id
  associate_public_ip_address = var.use_public_ip
  key_name                    = var.key_name
  vpc_security_group_ids      = [
    aws_security_group.awsbi_security_group.id
  ]
  tags = {
    Name = var.name
    cluster_name = var.name
  }
}

resource "aws_route_table" "awsbi_route_table" {
  vpc_id  = aws_vpc.awsbi_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.awsbi_internet_gateway.id
  }
  tags = {
        Name = "route-${var.name}"
        cluster_name = var.name
  }
}

resource "aws_route_table_association" "awsbi_route_association" {
  subnet_id      = aws_subnet.awsbi_subnet.id
  route_table_id = aws_route_table.awsbi_route_table.id
}
