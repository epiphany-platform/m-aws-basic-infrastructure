variable "name" {
  description = "Name to be used on all resources as prefix"
  type        = string
  default     = "awsbi"
}

variable "instance_type" {
  description = "The type of instance to start"
  type        = string
  default     = "t3.medium"
}

variable "instance_count" {
  description = "Number of instances to launch"
  type        = number
}

variable "root_volume_size" {
  description = "The size of the root volume in gibibytes (GiB)"
  type        = number
}

variable "use_public_ip" {
  description = "If true, the EC2 instance will have associated public IP address"
  type        = bool
}

variable "force_nat_gateway" {
  description = "If true, the NAT gateway will be forcefully deployed"
  type        = bool
}

variable "region" {
  description = "Region to launch in"
  type        = string
}

variable "key_name" {
  description = "Key pair name"
  type        = string
}

variable "vpc_cidr_block" {
  description = "The cidr block of the VPC"
  default     = "10.1.0.0/20"
}

variable "subnet_private_cidr_block" {
  description = "The cidr block of the private subnet"
  default     = "10.1.1.0/24"
}

variable "subnet_public_cidr_block" {
  description = "The cidr block of the public subnet"
  default     = "10.1.2.0/24"
}

variable "os" {
  description = "Operating System to launch"
  type = string
}
