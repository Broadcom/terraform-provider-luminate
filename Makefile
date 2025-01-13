GOCMD=go
GOBUILD=$(GOCMD) build $(GOARGS)
GOTEST=$(GOCMD) test $(GOARGS)
VERSION := $(shell cat VERSION)

MAKEFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR=$(dir $(MAKEFILE_PATH))

OUTPUT_DIR=build
BINARY_NAME=terraform-provider-luminate

GO111MODULE=on

all: linux darwin windows darwin_arm64 checksum sign

linux:
	mkdir -p release || true
	export GOOS=linux; export CGO_ENABLED=0; $(GOBUILD) -o $(OUTPUT_DIR)/linux/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)_$(VERSION)_linux.zip $(OUTPUT_DIR)/linux/*

darwin:
	mkdir -p release || true
	export GOOS=darwin; $(GOBUILD) -o $(OUTPUT_DIR)/darwin/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)_$(VERSION)_darwin.zip $(OUTPUT_DIR)/darwin/*

windows:
	mkdir -p release || true
	export GOOS=windows; $(GOBUILD) -o $(OUTPUT_DIR)/windows/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)_$(VERSION)_windows.zip $(OUTPUT_DIR)/windows/*

checksum:
	shasum -a 256 release/*.zip > terraform-provider-$(BINARY_NAME)_$(VERSION)_SHA256SUMS

testacc:
	export LUMINATE_API_ENDPOINT="${TERRAFORM_ACCEPTANCE_TENANT}" && \
	export LUMINATE_API_CLIENT_ID="${TERRAFORM_ACCEPTANCE_CLIENT_ID}" && \
	export LUMINATE_API_CLIENT_SECRET="${TERRAFORM_ACCEPTANCE_CLIENT_SECRET}" && \
	export TF_ACC=1 && \
	export TF_LOG="ERROR" && \
	export TEST_GROUP_NAME="tf-acceptance" && \
	export TEST_AWS_ACCOUNT_ID="957040371666" && \
	export TEST_AWS_INTEGRATION_NAME="terraform-test" && \
	export TEST_SSH_CLIENT_NAME="tf-at-ssh-client" && \
	export TEST_SSH_CLIENT_DESCRIPTION="a good description" && \
	export TEST_USERNAME="tf-user@terraformat.luminatesec.com" && \
	export TEST_USER_ID="f75f45b8-d10d-4aa6-9200-5c6d60110430" && \
	export TEST_USER_ID2="ed974d59-1941-4584-9336-2a9ed35043f2" && \
	export TEST_SITE_REGION="us-west1" && \
	export TEST_IDP_ID="0a0524a3-44ae-43b2-9d79-6cb018136b6e" && \
    $(GOTEST) -p 1 -v  ./...

testacc_serial:
	export LUMINATE_API_ENDPOINT="${LUMINATE_API_ENDPOINT}" && \
	export LUMINATE_API_CLIENT_ID="${LUMINATE_API_CLIENT_ID}" && \
	export LUMINATE_API_CLIENT_SECRET="${LUMINATE_API_CLIENT_SECRET}" && \
	export TF_ACC=1 && \
	export TF_LOG=ERROR && \
	export TEST_AWS_ACCOUNT_ID="${TEST_AWS_ACCOUNT_ID}" && \
	export TEST_AWS_INTEGRATION_NAME="${TEST_AWS_INTEGRATION_NAME}" && \
    $(GOTEST) -p 1 -v  ./provider/serial_tests

testacc_no_serial:
	export LUMINATE_API_ENDPOINT="${LUMINATE_API_ENDPOINT}" && \
	export LUMINATE_API_CLIENT_ID="${LUMINATE_API_CLIENT_ID}" && \
	export LUMINATE_API_CLIENT_SECRET="${LUMINATE_API_CLIENT_SECRET}" && \
	export TF_ACC=1 && \
	export TF_LOG=ERROR && \
	export TEST_GROUP_NAME="${TEST_GROUP_NAME}" && \
	export TEST_SSH_CLIENT_NAME="${TEST_SSH_CLIENT_NAME}" && \
	export TEST_SSH_CLIENT_DESCRIPTION="${TEST_SSH_CLIENT_DESCRIPTION}" && \
	export TEST_USERNAME="${TEST_USERNAME}" && \
	export TEST_USER_ID="${TEST_USER_ID}" && \
	export TEST_USER_ID2="${TEST_USER_ID2}" && \
	export TEST_SITE_REGION="${TEST_SITE_REGION}" && \
	export TEST_IDP_ID="${TEST_IDP_ID}" && \
	go_list_results=$$(go list ./... | grep -v 'serial_tests\|wss_tests') && $(GOTEST) -p 1 -v $$go_list_results

testacc_wss:
	export LUMINATE_API_ENDPOINT="${LUMINATE_API_ENDPOINT}" && \
	export LUMINATE_API_CLIENT_ID="${LUMINATE_API_CLIENT_ID}" && \
	export LUMINATE_API_CLIENT_SECRET="${LUMINATE_API_CLIENT_SECRET}" && \
	export TF_ACC=1 && \
	export TF_LOG=ERROR && \
	export  RUN_WSS_TESTS=true && \
    $(GOTEST) -p 1 -v  ./provider/wss_tests

darwin_arm64:
	mkdir -p release || true
	export GOOS=darwin GOARCH=arm64; $(GOBUILD) -o $(OUTPUT_DIR)/darwin_arm64/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-darwin_arm64.zip $(OUTPUT_DIR)/darwin_arm64/*

generate-docs:
	cd tools; go generate ./...

sign:
	gpg --detach-sign --output release/terraform-provider-$(NAME)_$(VERSION)_SHA256SUMS.sig release/terraform-provider-$(NAME)_$(VERSION)_SHA256SUMS
