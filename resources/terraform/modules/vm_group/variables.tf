variable "vpc_id" {
  description = "ID of VPC"
  type        = string
}

variable "name" {
  description = "String value to use as resources name prefix"
  type        = string
}

variable "key_name" {
  description = "Key pair name"
  type        = string
}

variable vm_group {
  description = "The list of VM group definition objects"
  type        = object({
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
  })

}
