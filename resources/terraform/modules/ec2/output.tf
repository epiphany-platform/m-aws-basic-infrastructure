output "private_ip" {
  value = aws_instance.awsbi.*.private_ip
}

output "public_ip" {
  value = aws_instance.awsbi.*.public_ip
}

output "vpc_id" {
  value = aws_vpc.awsbi_vpc.id
}

output "private_subnet_id" {
  value = aws_subnet.awsbi_private_subnet.id
}

output "public_subnet_id" {
  value = aws_subnet.awsbi_public_subnet.id
}
