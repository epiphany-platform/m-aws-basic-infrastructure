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

M_VMS_COUNT ?= 1
M_PUBLIC_IPS ?= false
M_NAT_GATEWAY_COUNT ?= 1
M_SUBNETS ?= $(_M_SUBNETS)
M_REGION ?= eu-central-1
M_NAME ?= epiphany
M_VMS_RSA ?= vms_rsa
M_OS ?= redhat
M_ADDRESS_SPACE ?= 10.1.0.0/20

AWS_ACCESS_KEY_ID ?= unset
AWS_SECRET_ACCESS_KEY ?= unset
