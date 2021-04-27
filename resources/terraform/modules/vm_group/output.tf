output "vm_group" {
  value = {
    vm_group_name: var.vm_group.name,
    vms: [
    for vm in aws_instance.awsbi:
    {
      vm_name: vm.tags.Name
      public_ip: vm.public_ip
      private_ip: vm.private_ip
      id: vm.id
    }
    ]
    data_disks: [
    for md in aws_ebs_volume.data_disks:
    {
      id: md.id
      name: md.tags.Name
      size: md.size
    }
    ]
    dd_attachments: [
    for dda in aws_volume_attachment.vms-dds-attachment:
    {
      data_disk_id: dda.volume_id
      instance_id: dda.instance_id
      device_name: dda.device_name
    }
    ]
  }
}
