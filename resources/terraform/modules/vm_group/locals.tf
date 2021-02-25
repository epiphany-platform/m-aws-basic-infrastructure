locals {
  vms_data_disks_product = setproduct(range(var.vm_group.vm_count), range(length(var.vm_group.data_disks)))
}
