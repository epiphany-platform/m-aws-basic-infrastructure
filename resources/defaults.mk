define _M_SUBNETS
{
  private:
  [{
    name: first_private_subnet,
    availability_zone: any,
    address_prefixes: 10.1.1.0/24
  }],
  public:
  [{
    name: first_public_subnet,
    availability_zone: any,
    address_prefixes: 10.1.2.0/24
  }]
}
endef

define _M_SECURITY_GROUPS
[
  {
    name: default_sg,
    rules: {
      ingress: [
        {
          protocol: "-1",
          from_port: 0,
          to_port: 0,
          cidr_blocks: ["10.1.0.0/20"]
        },
        {
          protocol: "tcp",
          from_port: 22,
          to_port: 22,
          cidr_blocks: ["0.0.0.0/0"]
        }
      ],
      egress: [
        {
          protocol: "-1",
          from_port: 0,
          to_port: 0,
          cidr_blocks: ["0.0.0.0/0"]
        }
      ]
    }
  }
]
endef

define _M_VM_GROUPS
[
  {
    name: vm-group0,
    vm_count: 1,
    vm_size: "t3.medium",
    use_public_ip: false,
    subnet_names: ["first_private_subnet"],
    sg_names: ["default_sg"],
    vm_image: {
      ami: "RHEL-7.8_HVM_GA-20200225-x86_64-1-Hourly2-GP2", # In case of Ubuntu: "ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20200611"
      owner: "309956199498", # In case of Ubuntu: "099720109477"
    },
    data_disks: [
      {
        device_name: "/dev/sdf", #/dev/sd[f-p]
        disk_size_gb: 16,
        type: "gp2"
      }
    ]
  }
]
endef

M_NAT_GATEWAY_COUNT ?= 1
M_SUBNETS ?= $(_M_SUBNETS)
M_SECURITY_GROUPS ?= $(_M_SECURITY_GROUPS)
M_VM_GROUPS ?= $(_M_VM_GROUPS)
M_REGION ?= eu-central-1
M_NAME ?= epiphany
M_VMS_RSA ?= vms_rsa
M_ADDRESS_SPACE ?= 10.1.0.0/20

AWS_ACCESS_KEY_ID ?= unset
AWS_SECRET_ACCESS_KEY ?= unset
