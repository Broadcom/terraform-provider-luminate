#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e
# Exit if any variable is not set
set -u
# set -x # Uncomment for full shell trace during debugging

echo "--- Starting Release Process (GPG Setup and Make) ---"

# Validate that required secret environment variables are present from CircleCI
if [ -z "${GPG_PRIVATE_KEY_CONTENT_B64:-}" ]; then
  echo "ERROR: GPG_PRIVATE_KEY_CONTENT_B64 CI secret is not set."
  exit 1
fi
if [ -z "${GPG_PASSPHRASE:-}" ]; then
  echo "ERROR: GPG_PASSPHRASE CI secret is not set."
  exit 1
fi
if [ -z "${GPG_SIGN_KEY_ID:-}" ]; then
  echo "Warning: GPG_SIGN_KEY_ID CI secret is not set. GPG may try to use a default key."
  # Depending on your GPG setup, you might want to exit 1 here if a specific key is always required
fi

# 1. Create a temporary directory for the GPG keyring
# Use a script-local variable for the temp home path
SCRIPT_TEMP_GPG_HOME="$(mktemp -d)"
echo "Setting up temporary GPG home: ${SCRIPT_TEMP_GPG_HOME}"

# Set up a trap to clean up the temporary directory when this script exits
trap "echo 'Cleaning up temp GPG home from do_release.sh: ${SCRIPT_TEMP_GPG_HOME}' && rm -rf \"${SCRIPT_TEMP_GPG_HOME}\"" EXIT
echo "Cleanup trap set for ${SCRIPT_TEMP_GPG_HOME} within do_release.sh"

# 2. Import the private key from the Base64 secret variable into the temporary home
echo "Importing private GPG key from Base64 secret variable..."
DECODED_KEY_DATA=""
if DECODED_KEY_DATA=$(printf "%s" "$GPG_PRIVATE_KEY_CONTENT_B64" | base64 --decode 2>/dev/null); then
  echo "Base64 decoding successful."
else
  echo "ERROR: base64 --decode command failed. Check GPG_PRIVATE_KEY_CONTENT_B64 secret and base64 utility in image."
  exit 1
fi

if echo "$DECODED_KEY_DATA" | gpg --batch --yes --import --homedir "${SCRIPT_TEMP_GPG_HOME}"; then
  echo "Private GPG key imported successfully into temporary keyring."
else
  echo "ERROR: GPG key import failed after base64 decode. The decoded data was not valid PGP."
  echo "Listing contents of ${SCRIPT_TEMP_GPG_HOME} after FAILED import attempt:"
  ls -laR "${SCRIPT_TEMP_GPG_HOME}"
  exit 1
fi

# 3. Export environment variables that Make will use.
# These are set in *this* shell environment, which will also run 'make'.
export GNUPGHOME="${SCRIPT_TEMP_GPG_HOME}" # GPG commands will use this
export GPG_HOME="${SCRIPT_TEMP_GPG_HOME}"  # For the Makefile to pick up and pass to sign_release.sh
if [ -n "$GPG_SIGN_KEY_ID" ]; then # GPG_SIGN_KEY_ID comes from CircleCI context
    export GPG_SIGN_KEY="$GPG_SIGN_KEY_ID" # For the Makefile to derive GPG_KEY_OPT
    echo "GPG_SIGN_KEY set to: (masked)"
else
    echo "Warning: GPG_SIGN_KEY_ID CI secret not found or empty."
fi
echo "GPG_HOME set to: ${GPG_HOME}"
echo "--- GPG Setup Complete ---"

# 4. Run the Full Release Build Pipeline (Make)
echo "Starting 'make release' pipeline..."
echo "Current GPG_HOME for make (from env): '${GPG_HOME:-<not_set_or_empty>}'"
echo "Current GPG_SIGN_KEY for make (from env): '${GPG_SIGN_KEY:-<not_set_or_empty>}'"
echo "Current GNUPGHOME for make (from env): '${GNUPGHOME:-<not_set_or_empty>}'"

# Pipe the GPG passphrase from the secret variable to Make's stdin
# The passphrase will be consumed by sign_release.sh (which is called by 'make release')
echo "$GPG_PASSPHRASE" | make release

MAKE_EXIT_CODE=$? # Capture exit code of make

echo "--- 'make release' Pipeline Finished with exit code ${MAKE_EXIT_CODE} ---"

# The trap will clean up SCRIPT_TEMP_GPG_HOME when this script exits.
exit ${MAKE_EXIT_CODE} # Exit with the make command's exit code
