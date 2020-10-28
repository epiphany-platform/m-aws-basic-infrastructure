resource "aws_key_pair" "kp" {
  key_name_prefix   = "kp-${var.name}"
  public_key = file(var.rsa_pub_path)
}

module "ec2" {
  source = "./modules/ec2"

  name              = var.name
  instance_count    = var.instance_count
  root_volume_size  = var.root_volume_size
  use_public_ip     = var.use_public_ip
  force_nat_gateway = var.force_nat_gateway
  region            = var.region
  key_name          = aws_key_pair.kp.key_name
  os                = var.os
}
