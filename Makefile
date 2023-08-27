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
	export TF_ACC=1 && $(GOTEST) -p 1 -v  ./...

darwin_arm64:
	mkdir -p release || true
	export GOOS=darwin GOARCH=arm64; $(GOBUILD) -o $(OUTPUT_DIR)/darwin_arm64/$(BINARY_NAME) -v
	zip -j release/$(BINARY_NAME)-darwin_arm64.zip $(OUTPUT_DIR)/darwin_arm64/*