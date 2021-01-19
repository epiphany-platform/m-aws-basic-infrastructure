# m-aws-basic-infrastructure

Epiphany Module: AWS Basic Infrastructure

AwsBI module is reponsible for providing basic cloud resources (eg. resource groups, virtual networks, subnets, virtual
machines etc.) which will be used by upcoming modules.

# Basic usage

## Requirements

Requirements are listed in a separate [document](docs/REQUIREMENTS.md).

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

  This command will create configuration file of AwsBI module in /tmp/shared/awsbi/awsbi-config.yml. You can investigate
  what is stored in that file. Available parameters are listed in the [inputs](docs/INPUTS.adoc) document.

* Plan and apply AwsBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest plan M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest apply M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx
  ```

  Running those commands should create a bunch of AWS resources (resource group, vpc, subnet, ec2 instances and so on).
  You can verify it in AWS Management Console.

* Passing complex structures as a parameter value

  It is possible that you need to pass a complex structure parameter value, such as `M_SUBNETS` or `M_SECURITY_GROUPS`,
  it can be achieved by 2 ways:

    1. Changing `defaults.mk` and building new Docker image with another default value
    2. Passing it as usual with a long string:

    ```shell
    docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:0.0.1 init \
      M_VMS_COUNT=2 \
      M_PUBLIC_IPS=true \
      M_NAME=epiphany-modules-awsbi \
      M_REGION=eu-west-3 \
      M_SECURITY_GROUPS='[ { name: default_sg, rules: { ingress: [ { protocol: "-1", from_port: 0, to_port: 0, cidr_blocks: ["10.1.0.0/20"] }, { protocol: "tcp", from_port: 22, to_port: 22, cidr_blocks: ["0.0.0.0/0"] }, { protocol: "tcp", from_port: 443, to_port: 443, cidr_blocks: ["0.0.0.0/0"] } ], egress: [ { protocol: "-1", from_port: 0, to_port: 0, cidr_blocks: ["0.0.0.0/0"] } ] } } ]'
    ```

* Destroy created infrastructure

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest plan-destroy M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsbi:latest destroy M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx
  ```

## Run module with provided example

### Prepare config file

Prepare your own variables in vars.mk file to use in the building process. Sample file (
examples/basic_flow/vars.mk.sample):

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

## Integration tests execution

Prior to run integration tests on for AWS module specify variables on OS where you want to run tests:

- AWS_ACCESS_KEY_ID - this is your access key
- AWS_SECRET_ACCESS_KEY - this is your secret
- AWSBI_IMAGE_TAG - this is full tag of docker image that you want to test e.g. "epiphanyplatform/awsbi:0.0.1"

and after that run shell command:

```shell
  make test
```

## Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AWS provider    | 3.7.0   | https://github.com/terraform-providers/terraform-provider-aws | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-aws/blob/master/LICENSE) |
| Make                      | 4.3     | https://www.gnu.org/software/make/                    | [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                        | 3.3.4   | https://github.com/mikefarah/yq/                      | [MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
| aws-sdk-go                | 1.15.77 | https://github.com/aws/aws-sdk-go/                    | [Apache License 2.0](https://github.com/aws/aws-sdk-go/blob/master/LICENSE.txt) | 