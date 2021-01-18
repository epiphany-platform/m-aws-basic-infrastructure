define _M_SUBNETS
{
  private: {
    count: 1
  },
  public: {
    count: 1
  }
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

M_VMS_COUNT ?= 1
M_PUBLIC_IPS ?= false
M_NAT_GATEWAY_COUNT ?= 1
M_SUBNETS ?= $(_M_SUBNETS)
M_SECURITY_GROUPS ?= $(_M_SECURITY_GROUPS)
M_REGION ?= eu-central-1
M_NAME ?= epiphany
M_VMS_RSA ?= vms_rsa
M_OS ?= redhat

AWS_ACCESS_KEY_ID ?= unset
AWS_SECRET_ACCESS_KEY ?= unset
