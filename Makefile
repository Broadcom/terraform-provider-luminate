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
	export LUMINATE_API_ENDPOINT="api.styrrbac.luminate-ci.com" && \
	export LUMINATE_API_CLIENT_ID="ae7f194fa945e284ffe3997b21121892" && \
	export LUMINATE_API_CLIENT_SECRET="20f5cccde80989cd1983c0f8a829c3e548be235fe5e3592979de5625c4845b40" && \
	export TF_ACC=1 && $(GOTEST) -p 1 -v  ./...
