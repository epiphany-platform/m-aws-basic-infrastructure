ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

AWSBI := epiphanyplatform/awsbi:latest

#makes it easier to replace the value
M_NAME ?= awsbi-module-example

VMS_RSA_FILENAME = vms_rsa

-include $(ROOT_DIR)/vars.mk

export

.PHONY: all clean metadata init plan apply output destroy-plan destroy setup

all: init plan apply

clean:
	@rm -rf $(ROOT_DIR)/shared

metadata:
	@docker run --rm \
		-t $(AWSBI) \
		metadata

init: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI) \
		init \
		--name=$(M_NAME) \
		--vms_rsa=$(VMS_RSA_FILENAME)

plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e ACCESS_KEY_ID=$$AWS_ACCESS_KEY \
		-e SECRET_ACCESS_KEY=$$AWS_SECRET_KEY \
		-t $(AWSBI) \
		plan

apply: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e ACCESS_KEY_ID=$$AWS_ACCESS_KEY \
		-e SECRET_ACCESS_KEY=$$AWS_SECRET_KEY \
		-t $(AWSBI) \
		apply

output: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-t $(AWSBI) \
		output

destroy-plan: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e ACCESS_KEY_ID=$$AWS_ACCESS_KEY \
		-e SECRET_ACCESS_KEY=$$AWS_SECRET_KEY \
		-t $(AWSBI) \
		plan \
		--destroy

destroy: setup
	@docker run --rm \
		-v $(ROOT_DIR)/shared:/shared \
		-e ACCESS_KEY_ID=$$AWS_ACCESS_KEY \
		-e SECRET_ACCESS_KEY=$$AWS_SECRET_KEY \
		-t $(AWSBI) \
		destroy

setup: $(ROOT_DIR)/shared/$(VMS_RSA_FILENAME)

$(ROOT_DIR)/shared/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
