GOCMD=go
GOBUILD=$(GOCMD) build $(GOARGS)
GOTEST=$(GOCMD) test $(GOARGS)

VERSION := $(shell cat VERSION)
VERSION_NO_V=$(patsubst v%,%,$(VERSION))

MAKEFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR=$(dir $(MAKEFILE_PATH))

OUTPUT_DIR=build
PROVIDER_NAME=terraform-provider-luminate
# --- Define Final Output File Paths (These names are required by the Registry) ---

# We still need these variables for use within the recipe commands
LINUX_AMD64_ZIP_PATH = $(RELEASE_DIR)/${PROVIDER_NAME}_$(VERSION)_linux_amd64.zip
DARWIN_AMD64_ZIP_PATH = $(RELEASE_DIR)/${PROVIDER_NAME}_$(VERSION)_darwin_amd64.zip
DARWIN_ARM64_ZIP_PATH = $(RELEASE_DIR)/${PROVIDER_NAME}_$(VERSION)_darwin_arm64.zip
WINDOWS_AMD64_ZIP_PATH = $(RELEASE_DIR)/${PROVIDER_NAME}_$(VERSION)_windows_amd64.zip

# Manifest file name: terraform-provider-{NAME}_{VERSION}_manifest.json
MANIFEST_FILE = $(RELEASE_DIR)/${PROVIDER_NAME}_$(VERSION)_manifest.json

# Checksums file name: terraform-provider-{NAME}_{VERSION}_SHA256SUMS
CHECKSUMS_FILE_PATH = $(RELEASE_DIR)/$(PROVIDER_NAME)_$(VERSION_NO_V)_SHA256SUMS

# Signature file name: terraform-provider-{NAME}_{VERSION}_SHA256SUMS.sig (Binary Signature)
SIGNATURE_FILE_PATH = $(RELEASE_DIR)/${PROVIDER_NAME}_$(VERSION)_SHA256SUMS.sig

# --- Define Temporary Binary Paths (Name inside the zip: terraform-provider-{NAME}_v{VERSION}) ---
LINUX_AMD64_BIN_TMP = $(BIN_TMP_DIR)/linux_amd64/${PROVIDER_NAME}_$(VERSION)
DARWIN_AMD64_BIN_TMP = $(BIN_TMP_DIR)/darwin_amd64/${PROVIDER_NAME}_$(VERSION)
DARWIN_ARM64_BIN_TMP = $(BIN_TMP_DIR)/darwin_arm64/${PROVIDER_NAME}_$(VERSION)
WINDOWS_AMD64_BIN_TMP = $(BIN_TMP_DIR)/windows_amd64/${PROVIDER_NAME}_$(VERSION).exe
# Output directory for final release artifacts (zip, json, sums, sig)
RELEASE_DIR=dist/$(VERSION)
# Temporary directory to build raw binaries before zipping
BIN_TMP_DIR=tmp/bin/$(VERSION)

build: clean linux darwin_amd64 darwin_arm64 windows manifest sign

release: clean linux darwin_amd64 darwin_arm64 windows manifest sign
	@echo "----------------------------------------------------"
	@echo " Release assets prepared in: $(RELEASE_DIR)"
	@echo " Files ready for hosting:"
	@ls -lh $(RELEASE_DIR)
	@echo "----------------------------------------------------"

# --- Clean Target ---
clean:
	@echo "--> Cleaning build artifacts..."
	@rm -rf $(RELEASE_DIR) $(BIN_TMP_DIR)

linux: $(LINUX_AMD64_BIN_TMP)
	@echo "--> Zipping Linux AMD64 binary ($(notdir $(LINUX_AMD64_ZIP_PATH)))..."
	@mkdir -p $(dir $(LINUX_AMD64_ZIP_PATH)) # Create directory for the zip file
	# cd to temp binary dir and zip the binary with the correct internal name to the release dir
	@cd $(dir $<) && zip $(abspath $(LINUX_AMD64_ZIP_PATH)) $(notdir $<)
	@rm -rf $(dir $<) # Clean up temp binary

$(LINUX_AMD64_BIN_TMP):
	@echo "--> Building Linux AMD64 binary ($(notdir $@))..."
	@mkdir -p $(dir $@) # Create directory for the temporary binary
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $@ . # <-- Build to the temporary path

darwin_amd64: $(DARWIN_AMD64_BIN_TMP)
	@echo "--> Zipping Darwin AMD64 binary ($(notdir $(DARWIN_AMD64_ZIP_PATH)))..."
	@mkdir -p $(dir $(DARWIN_AMD64_ZIP_PATH))
	@cd $(dir $<) && zip $(abspath $(DARWIN_AMD64_ZIP_PATH)) $(notdir $<)
	@rm -rf $(dir $<)

$(DARWIN_AMD64_BIN_TMP):
	@echo "--> Building Darwin AMD64 binary ($(notdir $@))..."
	@mkdir -p $(dir $@)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $@ .

darwin_arm64: $(DARWIN_ARM64_BIN_TMP)
	@echo "--> Zipping Darwin ARM64 binary ($(notdir $(DARWIN_ARM64_ZIP_PATH)))..."
	@mkdir -p $(dir $(DARWIN_ARM64_ZIP_PATH))
	@cd $(dir $<) && zip $(abspath $(DARWIN_ARM64_ZIP_PATH)) $(notdir $<)
	@rm -rf $(dir $<)

$(DARWIN_ARM64_BIN_TMP):
	@echo "--> Building Darwin ARM64 binary ($(notdir $@))..."
	@mkdir -p $(dir $@)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $@ .


windows: $(WINDOWS_AMD64_BIN_TMP)
	@echo "--> Zipping Windows AMD64 binary ($(notdir $(WINDOWS_AMD64_ZIP_PATH)))..."
	@mkdir -p $(dir $(WINDOWS_AMD64_ZIP_PATH))
	@cd $(dir $<) && zip $(abspath $(WINDOWS_AMD64_ZIP_PATH)) $(notdir $<)
	@rm -rf $(dir $<)

$(WINDOWS_AMD64_BIN_TMP):
	@echo "--> Building Windows AMD64 binary ($(notdir $@))..."
	@mkdir -p $(dir $@)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $@ .

manifest:
	@echo "--> Running script to create manifest ($(notdir $(PROVIDER_NAME)_$(VERSION)_manifest.json))..."
	@./.circleci/generate_manifest.sh $(RELEASE_DIR) $(PROVIDER_NAME) $(VERSION) # <--- This is the correct line

checksums:
	@echo "--> Generating SHA256 checksums ($(notdir $(CHECKSUMS_FILE_PATH)))..."
	# Ensure the release directory exists before writing files to it
	@mkdir -p $(RELEASE_DIR)
	# Navigate to the release directory, find files, checksum them, and write to the output file path
	# We write to the full $(CHECKSUMS_FILE_PATH) directly from the parent directory context
	find $(RELEASE_DIR) -type f \( -name "*.zip" -o -name "*.json" \) -print0 | sort -z | xargs -0 shasum -a 256 > $(CHECKSUMS_FILE_PATH) # Or sha256sum
	# Note: Removed 'cd $(RELEASE_DIR) &&' to avoid potential issues with relative paths in find/xargs and stdout redirection

	@echo "Generated checksums file:"
	@cat $(CHECKSUMS_FILE_PATH)

 # --- GPG Variables ---
 # Set GNUPGHOME if your GPG configuration/keys are not in the default location (~/.gnupg)
  GPG_HOME ?= $(HOME)/.gnupg

$(eval GPG_KEY_OPT = $(if $(GPG_SIGN_KEY),--local-user $(GPG_SIGN_KEY)))

sign: checksums $(CHECKSUMS_FILE_PATH)
	@echo "--> Generating binary GPG signature ($(notdir $(SIGNATURE_FILE_PATH)))..."
	set -x # Trace shell commands
	# Call the signing script, passing required arguments
	# Pass the passphrase to the script's stdin via pipe
	# Arguments: 1=CHECKSUMS_FILE_PATH, 2=SIGNATURE_FILE_PATH, 3=GPG_KEY_OPT, 4=GPG_HOME
	@./.circleci/sign_release.sh $(CHECKSUMS_FILE_PATH) $(SIGNATURE_FILE_PATH) "$(GPG_KEY_OPT)" "$(GPG_HOME)"
	@echo "Binary GPG signature generated."

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


generate-docs:
	cd tools; go generate ./...
