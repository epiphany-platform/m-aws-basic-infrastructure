ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

VERSION ?= dev
IMAGE_REPOSITORY := epiphanyplatform/awsbi

IMAGE_NAME := $(IMAGE_REPOSITORY):$(VERSION)

define AWS_CREDENTIALS_CONTENT
AWS_ACCESS_KEY_ID ?= $(ACCESS_KEY_ID)
AWS_SECRET_ACCESS_KEY ?= $(SECRET_ACCESS_KEY)
endef

-include ./aws-credentials.mk

export

#used for correctly setting shared folder permissions
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)

.PHONY: all

all: build

.PHONY: build test pipeline-test release prepare-aws-credentials

build: guard-IMAGE_NAME
	docker build \
		--progress plain \
		--build-arg ARG_M_VERSION=$(VERSION) \
		--build-arg ARG_HOST_UID=$(HOST_UID) \
		--build-arg ARG_HOST_GID=$(HOST_GID) \
		-t $(IMAGE_NAME) \
		.

#prepare AWS credentials variables file before running this target using `AWS_ACCESS_KEY_ID=xxx AWS_SECRET_ACCESS_KEY=yyy make prepare-aws-credentials`
test: guard-IMAGE_REPOSITORY build
	$(eval LDFLAGS = $(shell govvv -flags -pkg github.com/epiphany-platform/m-azure-basic-infrastructure/cmd -version $(VERSION)))
	@AWS_SECRET_ACCESS_KEY=$(AWS_ACCESS_KEY_ID) AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) AWSBI_IMAGE_TAG=$(IMAGE_REPOSITORY):$(VERSION) go test -ldflags="$(LDFLAGS)" -v -timeout 30m ./...

pipeline-test:
	$(eval LDFLAGS = $(shell govvv -flags -pkg github.com/epiphany-platform/m-azure-basic-infrastructure/cmd -version $(VERSION)))
	@go test -ldflags="$(LDFLAGS)" -v -timeout 30m

prepare-aws-credentials: guard-ACCESS_KEY_ID guard-SECRET_ACCESS_KEY
	@echo "$$AWS_CREDENTIALS_CONTENT" > $(ROOT_DIR)/aws-credentials.mk

release: guard-VERSION guard-IMAGE_NAME
	docker build \
		--build-arg ARG_M_VERSION=$(VERSION) \
		-t $(IMAGE_NAME) \
		.

print-%:
	@echo "$($*)"

guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

doctor:
	go mod tidy
	go fmt ./...
	go vet ./...
	goimports -l -w .
