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

variable "nat_gateway_count" {
  description = "The number of nat gateways to create"
  type        = number
}

variable "region" {
  description = "Region to launch in"
  type        = string
}

variable "key_name" {
  description = "Key pair name"
  type        = string
}

variable "vpc_address_space" {
  description = "The address space of the VPC"
  type        = string
}

variable "subnets" {
  type    = object({
    private = list(object({
      name                = string
      availability_zone   = string
      address_prefixes    = string
      }))
    public = list(object({
      name                = string
      availability_zone   = string
      address_prefixes    = string
      }))
  })
  validation {
    condition     = (length(var.subnets.private) > 0 && length(var.subnets.public) > 0) || length(var.subnets.public)> 0
    error_message = "Subnets list needs to have at least one element."
  }
}

variable "os" {
  description = "Operating System to launch"
  type = string
}
