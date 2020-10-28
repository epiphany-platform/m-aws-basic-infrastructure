include $(M_RESOURCES)/templates.mk
include $(M_RESOURCES)/consts.mk
include $(M_RESOURCES)/defaults.mk

export

.PHONY: metadata init plan apply audit destroy plan-destroy all-destroy output

#medatada method is printing static metadata information about module
metadata: guard-M_RESOURCES
	#AWSBI | metadata | Display component metadata
	@echo "$$M_METADATA_CONTENT"

#init method is used to initialize module configuration and check if state is providing strong (and weak) dependencies
#TODO should also validate state if strong requirements are met
init: guard-M_RESOURCES guard-M_SHARED guard-M_MODULE_SHORT guard-M_STATE_FILE_NAME \
			setup ensure-state-file template-config-file initialize-state-file display-config-file

#plan method would get config file and environment state file and compare them and calculate what would be done o apply stage
plan: guard-M_RESOURCES guard-M_SHARED guard-M_MODULE_SHORT guard-M_STATE_FILE_NAME \
			setup validate-config validate-state template-tfvars module-plan terraform-plan

#apply method runs module provider logic using config file
apply: guard-M_RESOURCES guard-M_SHARED \
			 setup module-plan terraform-apply update-state-after-apply terraform-output

#audit method checks if remote components are in "known" state
#TODO implement validation if remote resources are as expected, possibly with terraform plan
audit:
	#AWSBI | audit | should output current state of remote components

destroy: template-tfvars terraform-destroy update-state-after-destroy

plan-destroy: template-tfvars terraform-plan-destroy

all-destroy: plan-destroy destroy

output: terraform-output

setup: $(M_SHARED)/$(M_MODULE_SHORT)
	#AWSBI | setup | Ensure required directories exist

ensure-state-file: $(M_SHARED)/$(M_STATE_FILE_NAME)
	#AWSBI | ensure-state-file | Checks if 'state' file exists

template-config-file:
	#AWSBI | template-config-file | will template config file (and backup previous if exists)
	@if test -f $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME); then mv $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME).backup ; fi
	@echo "$$M_CONFIG_CONTENT" > $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME)

initialize-state-file:
	#AWSBI | initialize-state-file | will initialize state file
	@echo "$$M_STATE_INITIAL" > $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-state.tmp
	@yq m -i -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-state.tmp
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-state.tmp

#TODO: Get rid of error: '[Makefile:170: module-plan] Error 1 (ignored)'
module-plan:
	#AWSBI | module-plan | will perform module plan
	@yq m -x $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) > $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-future-state.tmp
	@yq w -i $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-future-state.tmp kind state
	@- yq compare $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-future-state.tmp
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-future-state.tmp

#TODO consider parsing terraform plan output
terraform-plan:
	#AWSBI | terraform-plan | will run plan
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform plan \
		-no-color \
		-input=false \
		-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		-out=$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan \
		$(M_RESOURCES)/terraform

terraform-plan-json:
	#AWSBI | terraform-plan-json | will show plan in json
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform show \
		-no-color \
		-json \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan

terraform-apply:
	#AWSBI | terraform-apply | will run terraform apply
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-apply.tfplan

terraform-plan-destroy:
	#AWSBI | terraform-plan-destroy | will prepare plan of destruction
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform plan \
		-destroy \
		-no-color \
		-input=false \
		-var-file=$(M_RESOURCES)/terraform/vars.tfvars.json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		-out=$(M_SHARED)/$(M_MODULE_SHORT)/terraform-destroy.tfplan \
		$(M_RESOURCES)/terraform

terraform-destroy:
	#AWSBI | terraform-destroy | will destroy using plan of destruction
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	TF_WARN_OUTPUT_ERRORS=1 \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform apply \
		-no-color \
		-input=false \
		-auto-approve \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate \
		$(M_SHARED)/$(M_MODULE_SHORT)/terraform-destroy.tfplan

terraform-output:
	#AWSBI | terraform-output | will prepare terraform output
	@cd $(M_RESOURCES)/terraform ; \
	TF_IN_AUTOMATION=true \
	AWS_ACCESS_KEY_ID=$(M_AWS_ACCESS_KEY) \
	AWS_SECRET_ACCESS_KEY=$(M_AWS_SECRET_KEY) \
		terraform output \
		-no-color \
		-json \
		-state=$(M_SHARED)/$(M_MODULE_SHORT)/terraform.tfstate > $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json
	@yq r --prettyPrint $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json | yq r - --printMode pv '*.value' > $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml
	@yq p -i $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml $(M_MODULE_SHORT).output
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.json $(M_SHARED)/$(M_MODULE_SHORT)/output.tmp.yml


display-config-file:
	#AWSBI | display-config-file | config file content is:
	@cat $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME)

update-state-after-apply:
	#AWSBI | update-state-after-apply | will update state file after apply
	@cp $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-config.tmp.yml
	@yq d -i $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-config.tmp.yml kind
	@yq m -x -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-config.tmp.yml
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status applied
	@rm $(M_SHARED)/$(M_MODULE_SHORT)/AWSBI-config.tmp.yml

update-state-after-destroy:
	#AWSBI | update-state-after-destroy | will clean state file after destroy
	@yq d -i $(M_SHARED)/$(M_STATE_FILE_NAME) '$(M_MODULE_SHORT)'
	@yq w -i $(M_SHARED)/$(M_STATE_FILE_NAME) $(M_MODULE_SHORT).status destroyed

#TODO check if there is state file
#TODO check if there is config
assert-init-completed:
	#AWSBI | assert-init-completed | will check if all initialization steps are completed

#TODO validate if config is correct
#TODO consider https://github.com/santhosh-tekuri/jsonschema as it's small
validate-config:
	#AWSBI | validate-config | will perform config validation

#TODO validate if state file is correct
#TODO consider https://github.com/santhosh-tekuri/jsonschema as it's small
validate-state:
	#AWSBI | validate-state | will perform state file validation

template-tfvars:
	#AWSBI | template-tfvars | will template .tfvars.json file
	@yq read -jP $(M_SHARED)/$(M_MODULE_SHORT)/$(M_CONFIG_NAME) '$(M_MODULE_SHORT)*' > $(M_RESOURCES)/terraform/vars.tfvars.json


guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

$(M_SHARED)/$(M_STATE_FILE_NAME):
	@mkdir -p $(dir $@)
	@touch $@

$(M_SHARED)/$(M_MODULE_SHORT):
	@mkdir -p $@
