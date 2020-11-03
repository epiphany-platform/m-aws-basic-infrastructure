# m-aws-basic-infrastructure

Epiphany Module: AWS Basic Infrastructure

AwsBI module is reponsible for providing basic cloud resources (eg. resource groups, virtual networks, subnets, virtual machines etc.) which will be used by upcoming modules.

# Basic usage

## Build image

In main directory run:

  ```shell
  make build
  ```

or directly using Docker:

  ```shell
  cd m-aws-basic-infrastructure/
  docker build --tag epiphanyplatform/awsbi:latest .
  ```

## Run module

* Create a shared directory:

  ```shell
  mkdir /tmp/shared
  ```

  This 'shared' dir is a place where all configs and states will be stored while working with modules.

* Generate ssh keys in: /tmp/shared/vms_rsa.pub

  ```shell
  ssh-keygen -t rsa -b 4096 -f /tmp/shared/vms_rsa -N ''
  ```

* Initialize AwsBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest init M_VMS_COUNT=2 M_PUBLIC_IPS=true M_NAME=epiphany-modules-awsbi
  ```

  This command will create configuration file of AwsBI module in /tmp/shared/awsbi/awsbi-config.yml. You can investigate what is stored in that file.
  Available variables:
  * M_VMS_COUNT = < instance count > (default: 1)
  * M_PUBLIC_IPS =  < true/false > (default: true)
  * M_NAME = < module name > (default: epiphany)
  * M_VMS_RSA = < ssh priv key name > (defailt: vms_rsa)
  * M_REGION = < aws region > (default: eu-central-1)
  * M_OS = < ubuntu/redhat > (default: ubuntu)

* Plan and apply AwsBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest plan M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx M_NAME=xxx
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest apply M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx M_NAME=xxx
  ```

  Running those commands should create a bunch of AWS resources (resource group, vpc, subnet, ec2 instances and so on). You can verify it in AWS Management Console.

## Run module with provided example

### Prepare config file

Prepare your own variables in vars.mk file to use in the building process.
Sample file (examples/basic_flow/vars.mk.sample):

  ```shell
  AWS_ACCESS_KEY_ID = "xxx"
  AWS_ACCESS_KEY_SECRET = "xxx"
  ```

### Create an environment

  ```shell
  cd examples/basic_flow
  make all
  ```

or step-by-step:

  ```shell
  cd examples/basic_flow
  make init
  make plan
  make apply
  ```

### Delete environment

  ```shell
  cd examples/basic_flow
  make all-destroy
  ```

or step-by-step

  ```shell
  cd examples/basic_flow
  make destroy-plan
  make destroy
  ```

## Release module

  ```shell
  make release
  ```

or if you want to set a different version number:

  ```shell
  make release VERSION=number_of_your_choice
  ```

# Awsbi output data

The output from this module is:

* private_ip
* public_ip
* public_subnet_id
* vpc_id
* private_route_table_id

# Windows users

This module is designed for Linux/Unix development/usage only. If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then ```examples/basic_flow_devcontainer``` folder.

## Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AWS provider    | 3.7.0   | https://github.com/terraform-providers/terraform-provider-aws | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-aws/blob/master/LICENSE) |
| Make                      | 4.3     | https://www.gnu.org/software/make/                    | [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                        | 3.3.4   | https://github.com/mikefarah/yq/                      | [MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
