variable "name" {
  description = "Name to be used on all resources as prefix"
  type        = string
}

variable "region" {
  description = "Region to launch in"
  type        = string
}

variable "nat_gateway_count" {
  description = "The number of nat gateways to create"
  type        = number
}

variable "vpc_address_space" {
  description = "The address space of the VPC"
  type        = string
}

variable "subnets" {
  description = "Subnets configuration"
  type        = object({
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
    condition     = length(var.subnets.private) > 0 || length(var.subnets.public) > 0
    error_message = "You must specify at least one subnet."
  }
}

variable "security_groups" {
  description = "Security groups configuration"
  type        = list(object({
    name        = string
    rules       = object({
      ingress = list(object({
        protocol    = string
        from_port   = number
        to_port     = number
        cidr_blocks = list(string)
      }))
      egress = list(object({
        protocol    = string
        from_port   = number
        to_port     = number
        cidr_blocks = list(string)
      }))
    })
  }))
}

variable "rsa_pub_path" {
  type = string
}

variable vm_groups {
  description = "The list of VM group definition objects"
  type        = list(object({
    name          = string
    vm_count      = number
    vm_size       = string
    use_public_ip = bool
    subnet_names  = list(string)
    sg_names      = list(string)
    vm_image      = object({
      ami         = string
      owner       = string
    })
    data_disks    = list(object({
      device_name = string
      disk_size_gb = number
      type        = string
    }))
  }))
}
