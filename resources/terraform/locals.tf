locals {
  use_nat_gateway           = (length(var.subnets.public) > 0 && length(var.subnets.private) > 0 && var.nat_gateway_count > 0)
  use_internet_gateway      = length(var.subnets.public) > 0
}
