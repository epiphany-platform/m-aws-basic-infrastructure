variable "name" {
  description = "Name to be used on all resources as prefix"
  type        = string
}

variable "instance_count" {
  description = "Number of instances to launch"
  type        = number
}

variable "root_volume_size" {
  description = "The size of the root volume in gibibytes (GiB)"
  type        = number
  default     = 64
}

variable "region" {
  description = "Region to launch in"
  type        = string
}

variable "use_public_ip" {
  description = "If true, the EC2 instance will have associated public IP address"
  type        = bool
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
}

variable "rsa_pub_path" {
  type = string
}

variable "os" {
  description = "Operating System to launch"
  type = string
}
