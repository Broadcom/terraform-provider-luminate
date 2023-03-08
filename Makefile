GOCMD=go
#GOARGS="-mod=vendor"
GOBUILD=$(GOCMD) build $(GOARGS)
GOTEST=$(GOCMD) test $(GOARGS)

MAKEFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR=$(dir $(MAKEFILE_PATH))

OUTPUT_DIR=build
BINARY_NAME=terraform-provider-luminate

GO111MODULE=on

all: linux darwin windows

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
	export LUMINATE_API_ENDPOINT="api.****.luminate-ci.com" && \
	export LUMINATE_API_CLIENT_ID="20810a69a650b3562987576cc3bbb45f" && \
	export LUMINATE_API_CLIENT_SECRET="e9927434300e0dc51259492f638901e22db1a20bebb71bba8075cc5fe49f962a" && \
	export TF_ACC=1 && $(GOTEST) -p 1 -v  ./...
