data "aws_subnet" "selected" {
  count                = length(var.vm_group.subnet_names)
  vpc_id               = var.vpc_id
  filter {
    name   = "tag:Name"
    values = ["${var.vm_group.subnet_names[count.index]}"]
  }
}

data "aws_security_group" "selected" {
  count  = length(var.vm_group.sg_names)
  name   = var.vm_group.sg_names[count.index]
  vpc_id = var.vpc_id
}

data "aws_ami" "select" {
  owners = [var.vm_group.vm_image.owner]
  filter {
    name   = "name"
    values = ["${var.vm_group.vm_image.ami}"]
  }
  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "awsbi" {
  count                       = var.vm_group.vm_count
  ami                         = data.aws_ami.select.id
  instance_type               = var.vm_group.vm_size
  subnet_id                   = element(data.aws_subnet.selected.*.id, count.index)
  associate_public_ip_address = var.vm_group.use_public_ip
  key_name                    = var.key_name

  vpc_security_group_ids      = data.aws_security_group.selected.*.id 

  tags = {
    Name = "${var.name}-${var.vm_group.name}-${count.index}"
    resource_group = "${var.name}"
  }
}

resource "aws_ebs_volume" "data_disks" {
  count             = length(local.vms_data_disks_product)
  availability_zone = aws_instance.awsbi[local.vms_data_disks_product[count.index][0]].availability_zone
  type              = var.vm_group.data_disks[local.vms_data_disks_product[count.index][1]].type
  size              = var.vm_group.data_disks[local.vms_data_disks_product[count.index][1]].disk_size_gb
  
  tags = {
    Name = "${var.name}-${var.vm_group.name}-${local.vms_data_disks_product[count.index][0]}-data-disk-${local.vms_data_disks_product[count.index][1]}"
  }
}

resource "aws_volume_attachment" "vms-dds-attachment" {
  count       = length(local.vms_data_disks_product)
  device_name = var.vm_group.data_disks[local.vms_data_disks_product[count.index][1]].device_name
  instance_id = aws_instance.awsbi[local.vms_data_disks_product[count.index][0]].id
  volume_id   = aws_ebs_volume.data_disks[count.index].id
}
