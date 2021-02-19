resource "aws_key_pair" "kp" {
  key_name_prefix   = "${var.name}-kp"
  public_key = file(var.rsa_pub_path)
  tags = {
    resource_group = "${var.name}-rg"
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
      "Values": ["${var.name}-rg"]
    }
  ]
}
JSON
  }
}

module "vm_group" {
  source            = "./modules/vm_group"
  count             = length(var.vm_groups)
  vm_group          = var.vm_groups[count.index]
  key_name          = aws_key_pair.kp.key_name
  vpc_id            = aws_vpc.awsbi_vpc.id
  name              = var.name

  depends_on = [aws_subnet.awsbi_public_subnet]
}
