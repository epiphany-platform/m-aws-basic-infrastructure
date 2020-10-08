resource "aws_key_pair" "kp" {
  key_name   = "kp-${var.name}"
  public_key = file(var.rsa_pub_path)
}

module "ec2" {
  source = "./modules/ec2"

  name           = var.name
  instance_count = var.instance_count
  use_public_ip  = var.use_public_ip
  region         = var.region
  ami            = var.ami
  key_name       = aws_key_pair.kp.key_name
}
