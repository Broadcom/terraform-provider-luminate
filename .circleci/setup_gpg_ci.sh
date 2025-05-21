#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e
# Exit if any variable is not set
set -u
# set -x # Uncomment for full shell trace during debugging

echo "--- GPG Setup Script Start ---"

# Validate that required secret environment variables are present
if [ -z "${GPG_PRIVATE_KEY_CONTENT_B64:-}" ]; then
  echo "ERROR: GPG_PRIVATE_KEY_CONTENT_B64 secret is not set."
  exit 1
fi
if [ -z "${GPG_PASSPHRASE:-}" ]; then # Passphrase isn't used directly here, but good to check if it exists for the 'make' step
  echo "Warning: GPG_PASSPHRASE secret is not set. The 'make release' command will likely fail at signing."
  # Depending on your workflow, you might want to exit 1 here too.
fi
if [ -z "${GPG_SIGN_KEY_ID:-}" ]; then
  echo "Warning: GPG_SIGN_KEY_ID secret is not set. GPG may use default key if available in imported keyring."
fi

# 1. Create a temporary directory for the GPG keyring
# Use a script-local variable for the temp home path
CI_GPG_TEMP_HOME="$(mktemp -d)"
echo "Setting up temporary GPG home: ${CI_GPG_TEMP_HOME}"

# Set up a trap to clean up the temporary directory when this script exits
# Note: This trap is for the lifetime of THIS script.
# The main job trap in config.yml is still good practice for overall job cleanup.
trap "echo 'Cleaning up temp GPG home from setup_gpg_ci.sh: ${CI_GPG_TEMP_HOME}' && rm -rf \"${CI_GPG_TEMP_HOME}\"" EXIT
echo "Cleanup trap set for ${CI_GPG_TEMP_HOME} within setup_gpg_ci.sh"

# 2. Import the private key from the Base64 secret variable into the temporary home
echo "Importing private GPG key from Base64 secret variable..."
DECODED_KEY_DATA=""
if DECODED_KEY_DATA=$(printf "%s" "$GPG_PRIVATE_KEY_CONTENT_B64" | base64 --decode 2>/dev/null); then
  echo "Base64 decoding successful."
else
  echo "ERROR: base64 --decode command failed. Check GPG_PRIVATE_KEY_CONTENT_B64 secret and base64 utility in image."
  exit 1
fi

if echo "$DECODED_KEY_DATA" | gpg --batch --yes --import --homedir "${CI_GPG_TEMP_HOME}"; then
  echo "Private GPG key imported successfully into temporary keyring."
else
  echo "ERROR: GPG key import failed after base64 decode. The decoded data was not valid PGP."
  exit 1
fi

# Optional debug: Verify key listing in the temporary keyring
# echo "Listing secret keys in temporary keyring:"
# gpg --list-secret-keys --keyid-format LONG --homedir "${CI_GPG_TEMP_HOME}" || echo "Warning: Failed to list secret keys."

# 3. Export environment variables that Make will use
# These need to be available to the *parent shell* that calls `make`.
# The best way to do this from a script called by a 'run' step is to write them to $BASH_ENV.
echo "Exporting GPG environment variables to $BASH_ENV for subsequent steps..."
echo "export GNUPGHOME=\"${CI_GPG_TEMP_HOME}\"" >> $BASH_ENV
echo "export GPG_HOME=\"${CI_GPG_TEMP_HOME}\"" >> $BASH_ENV
if [ -n "$GPG_SIGN_KEY_ID" ]; then
    echo "export GPG_SIGN_KEY=\"$GPG_SIGN_KEY_ID\"" >> $BASH_ENV
    echo "GPG_SIGN_KEY will be set to (masked) via $BASH_ENV"
else
    echo "Warning: GPG_SIGN_KEY_ID CI secret not found or empty. GPG_SIGN_KEY will not be set via $BASH_ENV."
fi

echo "GPG_HOME will be set to ${CI_GPG_TEMP_HOME} via $BASH_ENV"
echo "--- GPG Setup Script Finished ---"
