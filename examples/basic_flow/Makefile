ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

AWSBI_IMAGE_NAME := epiphanyplatform/awsbi:0.0.1

-include $(ROOT_DIR)/vars.mk

M_NAME ?= awsbi-basic-flow
M_VMS_RSA_FILENAME ?= vms_rsa

export

.PHONY: all metadata setup init plan apply output all-destroy plan-destroy destroy clean

warning:
	$(error Usage: make (all/metadata/setup/init/plan/apply/output/all-destroy/destroy/destroy-plan/clean) )

all: init plan apply
all-destroy: plan-destroy destroy

clean:
	@rm -rf $(ROOT_DIR)/shared

metadata: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		metadata \

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		init \
		M_NAME=$(M_NAME)

plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		plan \
		M_AWS_ACCESS_KEY=$(AWS_ACCESS_KEY) \
		M_AWS_SECRET_KEY=$(AWS_SECRET_KEY)

apply: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		apply \
		M_AWS_ACCESS_KEY=$(AWS_ACCESS_KEY) \
		M_AWS_SECRET_KEY=$(AWS_SECRET_KEY)

output: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		output \
		M_AWS_ACCESS_KEY=$(AWS_ACCESS_KEY) \
		M_AWS_SECRET_KEY=$(AWS_SECRET_KEY)

plan-destroy: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		plan-destroy \
		M_AWS_ACCESS_KEY=$(AWS_ACCESS_KEY) \
		M_AWS_SECRET_KEY=$(AWS_SECRET_KEY)

destroy: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI_IMAGE_NAME) \
		destroy \
		M_AWS_ACCESS_KEY=$(AWS_ACCESS_KEY) \
		M_AWS_SECRET_KEY=$(AWS_SECRET_KEY)

setup: $(ROOT_DIR)/shared/$(M_VMS_RSA_FILENAME)

$(ROOT_DIR)/shared/$(M_VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
