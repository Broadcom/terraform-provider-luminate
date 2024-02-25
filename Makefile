GOCMD=go
GOBUILD=$(GOCMD) build $(GOARGS)
GOTEST=$(GOCMD) test $(GOARGS)

MAKEFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR=$(dir $(MAKEFILE_PATH))

OUTPUT_DIR=build
BINARY_NAME=terraform-provider-luminate

GO111MODULE=on

all: linux darwin windows darwin_arm64

linux:
	mkdir -p release || true
	export GOOS=linux; $(GOBUILD) -o $(OUTPUT_DIR)/linux/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-linux.zip $(OUTPUT_DIR)/linux/*

darwin:
	mkdir -p release || true
	export GOOS=darwin; $(GOBUILD) -o $(OUTPUT_DIR)/darwin/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-darwin.zip $(OUTPUT_DIR)/darwin/*

windows:
	mkdir -p release || true
	export GOOS=windows; $(GOBUILD) -o $(OUTPUT_DIR)/windows/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-windows.zip $(OUTPUT_DIR)/windows/*

testacc:
	export LUMINATE_API_ENDPOINT="${TERRAFORM_ACCEPTANCE_TENANT}" && \
	export LUMINATE_API_CLIENT_ID="${TERRAFORM_ACCEPTANCE_CLIENT_ID}" && \
	export LUMINATE_API_CLIENT_SECRET="${TERRAFORM_ACCEPTANCE_CLIENT_SECRET}" && \
	export TF_ACC=1 && \
	export TF_LOG=DEBUG && \
	export TEST_GROUP_NAME="${TEST_GROUP_NAME}" && \
	export TEST_AWS_ACCOUNT_ID="${TEST_AWS_ACCOUNT_ID}" && \
	export TEST_AWS_INTEGRATION_NAME="${TEST_AWS_INTEGRATION_NAME}" && \
	export TEST_SSH_CLIENT_NAME="${TEST_SSH_CLIENT_NAME}" && \
	export TEST_SSH_CLIENT_ID="${TEST_SSH_CLIENT_ID}" && \
	export TEST_USERNAME="${TEST_USERNAME}" && \
	export TEST_USER_ID="${TEST_USER_ID}" && \
	export TEST_USER_ID2="${TEST_USER_ID2}" && \
	export TEST_SITE_REGION="${TEST_SITE_REGION}" && \
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
	export TEST_SSH_CLIENT_ID="${TEST_SSH_CLIENT_ID}" && \
	export TEST_USERNAME="${TEST_USERNAME}" && \
	export TEST_USER_ID="${TEST_USER_ID}" && \
	export TEST_USER_ID2="${TEST_USER_ID2}" && \
	export TEST_SITE_REGION="${TEST_SITE_REGION}" && \
	$(GOTEST) -p 1 -v `go list ./... | grep -v 'serial_tests'`

darwin_arm64:
	mkdir -p release || true
	export GOOS=darwin GOARCH=arm64; $(GOBUILD) -o $(OUTPUT_DIR)/darwin_arm64/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-darwin_arm64.zip $(OUTPUT_DIR)/darwin_arm64/*