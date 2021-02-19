define M_METADATA_CONTENT
labels:
  version: $(M_VERSION)
  name: AWS Basic Infrastructure
  short: $(M_MODULE_SHORT)
  kind: infrastructure
  provider: aws
  provides-vms: true
  provides-pubips: $(M_PUBLIC_IPS)
endef

define M_CONFIG_CONTENT
kind: $(M_MODULE_SHORT)-config
$(M_MODULE_SHORT):
  name: $(M_NAME)
  region: $(M_REGION)
  nat_gateway_count: $(M_NAT_GATEWAY_COUNT)
  rsa_pub_path: "$(M_SHARED)/$(M_VMS_RSA).pub"
  vpc_address_space: $(M_ADDRESS_SPACE)  
  subnets: $(M_SUBNETS)
  security_groups: $(M_SECURITY_GROUPS)
  vm_groups: $(M_VM_GROUPS)
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
