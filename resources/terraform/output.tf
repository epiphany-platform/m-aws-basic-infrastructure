output "vpc_id" {
  value = aws_vpc.awsbi_vpc.id
}

output "private_subnet_ids" {
  value = aws_subnet.awsbi_private_subnet.*.id
}

output "public_subnet_ids" {
  value = aws_subnet.awsbi_public_subnet.*.id
}

output "private_route_table" {
  value = join(" ", aws_route_table.awsbi_route_table_private.*.id)
}

output "vm_group" {
  value = module.vm_group.*.vm_group
}
