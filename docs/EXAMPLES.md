# Examples

## Prepare docker image

```shell
make build VERSION=latest
```

## Prepare credentials

Prepare your own variables in `vars.mk` file to use in the process.
Sample file (`examples/basic_flow/vars.mk.sample`):

```shell
AWS_ACCESS_KEY ?= aws_access_key_id
AWS_SECRET_KEY ?= aws_secret_access_key
```

# Create cluster

```shell
cd examples/basic_flow
make all
```

# Delete cluster

```shell
make destroy-plan
make destroy
```
