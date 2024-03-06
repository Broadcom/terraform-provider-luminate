GOCMD=go
GOBUILD=$(GOCMD) build $(GOARGS)
GOTEST=$(GOCMD) test $(GOARGS)

MAKEFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR=$(dir $(MAKEFILE_PATH))

OUTPUT_DIR=build
BINARY_NAME=terraform-provider-luminate

GO111MODULE=on
define TEST_ENV_VARS
export LUMINATE_API_ENDPOINT="${LUMINATE_API_ENDPOINT}"
export LUMINATE_API_CLIENT_ID="${LUMINATE_API_CLIENT_ID}"
export LUMINATE_API_CLIENT_SECRET="${LUMINATE_API_CLIENT_SECRET}"
export TF_ACC=1
export TF_LOG=ERROR
export TEST_GROUP_NAME="${TEST_GROUP_NAME}"
export TEST_SSH_CLIENT_NAME="${TEST_SSH_CLIENT_NAME}"
export TEST_SSH_CLIENT_ID="${TEST_SSH_CLIENT_ID}"
export TEST_USERNAME="${TEST_USERNAME}"
export TEST_USER_ID="${TEST_USER_ID}"
export TEST_USER_ID2="${TEST_USER_ID2}"
export TEST_SITE_REGION="${TEST_SITE_REGION}"
endef

testacc_no_serial:
	$(TEST_ENV_VARS)
	@echo "Running go list command"
	go_list_results=`go list ./... | grep -v 'serial_tests'`
	@echo "go list command executed successfully"
	@echo "Running go test command"
	$(GOTEST) -p 1 -v go_list_results
	@echo "go test command executed successfully."

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
	export LUMINATE_API_ENDPOINT="${api.terraformat.luminatesec.com}" && \
	export LUMINATE_API_CLIENT_ID="${16a49c4909f4165f7ad33f1d022d83fc}" && \
	export LUMINATE_API_CLIENT_SECRET="${734ea9a0d4ae44c55f00773a14538f505d16058fc892286b17a0219be4f84c54}" && \
	export TF_ACC=1 && \
	export TF_LOG=ERROR && \
	export TEST_GROUP_NAME="${tf-acceptance}" && \
	export TEST_AWS_ACCOUNT_ID="${957040371666}" && \
	export TEST_AWS_INTEGRATION_NAME="${terraform-test}" && \
	export TEST_SSH_CLIENT_NAME="${tf-at-ssh-client}" && \
	export TEST_SSH_CLIENT_ID="${6ddace8e-39a3-4cd3-bba7-ad26e826df5b}" && \
	export TEST_USERNAME="${tf-user@terraformat.luminatesec.com}" && \
	export TEST_USER_ID="${f75f45b8-d10d-4aa6-9200-5c6d60110430}" && \
	export TEST_USER_ID2="${ed974d59-1941-4584-9336-2a9ed35043f2}" && \
	export TEST_SITE_REGION="${us-west1}" && \
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
	go_list_results=$$(go list ./... | grep -v 'serial_tests') && $(GOTEST) -p 1 -v $$go_list_results

darwin_arm64:
	mkdir -p release || true
	export GOOS=darwin GOARCH=arm64; $(GOBUILD) -o $(OUTPUT_DIR)/darwin_arm64/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-darwin_arm64.zip $(OUTPUT_DIR)/darwin_arm64/*