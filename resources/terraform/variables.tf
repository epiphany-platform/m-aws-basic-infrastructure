variable "name" {
  description = "Name to be used on all resources as prefix"
  type        = string
}

variable "instance_count" {
  description = "Number of instances to launch"
  type        = number
}

variable "region" {
  description = "Region to launch in"
  type        = string
}

variable "use_public_ip" {
  description = "If true, the EC2 instance will have associated public IP address"
  type        = bool
}

variable "force_nat_gateway" {
  description = "If true, the NAT gateway will be forcefully deployed"
  type        = bool
}

variable "rsa_pub_path" {
  type = string
}

variable "admin_username" {
  type    = string
  default = "ec2-user"
}

variable "ami" {
  description = "ID of AMI to use for the instance"
  type        = string
  default     = "ami-02ab606eae7264892"
}
