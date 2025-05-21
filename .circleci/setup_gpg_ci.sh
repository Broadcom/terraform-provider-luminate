#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e
# Exit if any variable is not set
set -u
# set -x # Uncomment for full shell trace during debugging

echo "--- GPG Setup Script Start ---"

# Validate that required secret environment variables are present
# ... (your existing validation for GPG_PRIVATE_KEY_CONTENT_B64, etc.) ...

# 1. Create a temporary directory for the GPG keyring
CI_GPG_TEMP_HOME="$(mktemp -d)"
echo "Setting up temporary GPG home: ${CI_GPG_TEMP_HOME}"
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
  # Add ls here too if import fails, to see what was (or wasn't) created
  echo "Listing contents of ${CI_GPG_TEMP_HOME} after FAILED import attempt:"
  ls -laR "${CI_GPG_TEMP_HOME}"
  exit 1
fi

# --- ADD DEBUGGING HERE ---
echo "Listing secret keys in temporary keyring after import:"
gpg --list-secret-keys --keyid-format LONG --homedir "${CI_GPG_TEMP_HOME}" || echo "Warning: Failed to list secret keys after import."
echo "Listing contents of ${CI_GPG_TEMP_HOME} after successful import:"
ls -laR "${CI_GPG_TEMP_HOME}" # <--- Add this line
# --- END DEBUGGING ---


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
