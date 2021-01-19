data "aws_ami" "select" {
  owners = [local.select_owner]
  filter {
    name   = "name"
    values = ["${local.select_ami}"]
  }
  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_resourcegroups_group" "rg" {
  name = "${var.name}-rg"

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
    "AWS::EC2::RouteTable",
    "AWS::EC2::EIP"
  ],
  "TagFilters": [
    {
      "Key": "resource_group",
      "Values": ["${var.name}"]
    }
  ]
}
JSON
  }
}

resource "aws_instance" "awsbi" {
  count                       = var.instance_count
  ami                         = data.aws_ami.select.id
  instance_type               = var.instance_type
  subnet_id                   = var.use_public_ip ? element(aws_subnet.awsbi_public_subnet.*.id, count.index) : element(aws_subnet.awsbi_private_subnet.*.id, count.index)
  associate_public_ip_address = var.use_public_ip
  key_name                    = var.key_name

  root_block_device {
    volume_size = var.root_volume_size
  }

  vpc_security_group_ids = aws_security_group.awsbi_security_group.*.id

  tags = {
    Name = "${var.name}-instance${count.index}"
    resource_group = var.name
  }
}
