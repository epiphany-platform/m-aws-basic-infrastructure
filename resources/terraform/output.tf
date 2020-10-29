output "private_ip" {
  value = module.ec2.private_ip
}

output "public_ip" {
  value = module.ec2.public_ip
}

output "vpc_id" {
  value = module.ec2.vpc_id
}

output "public_subnet_id" {
  value = module.ec2.public_subnet_id
}

output "private_route_table_id" {
  value = module.ec2.private_route_table
}
