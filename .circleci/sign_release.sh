#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e
# Exit if any variable is not set
set -u

# set -x # Uncomment for full shell trace

# --- Script Inputs ---
# Pass these as command-line arguments from the Makefile
CHECKSUMS_FILE_REL="$1"   # Path to the SHA256SUMS file (relative)
SIGNATURE_FILE_REL="$2"   # Path where the binary signature should be saved (relative)
GPG_KEY_OPT="$3"          # Optional: GPG '--local-user KEY_ID' string
GPG_HOME_PATH="$4"        # Optional: Path for GNUPGHOME

# Convert relative paths to absolute paths
# This assumes the script is called from the project root directory
CHECKSUMS_FILE="$(pwd)/${CHECKSUMS_FILE_REL}"
SIGNATURE_FILE="$(pwd)/${SIGNATURE_FILE_REL}"


# --- Validation ---
# Validate original arguments for emptiness
if [ -z "$CHECKSUMS_FILE_REL" ] || [ -z "$SIGNATURE_FILE_REL" ]; then
  echo "Usage: $0 <checksums_file_rel> <signature_file_rel> [gpg_key_opt] [gnupghome_path]"
  exit 1
fi

# Validate that the absolute checksums file exists
if [ ! -f "$CHECKSUMS_FILE" ]; then
  echo "Error: Checksums file not found: ${CHECKSUMS_FILE}"
  exit 1
fi

# Ensure the output directory for the signature exists (using absolute path)
mkdir -p "$(dirname "${SIGNATURE_FILE}")"


echo "--> Running signing script..."
echo "    Signing file (relative input): ${CHECKSUMS_FILE_REL}"
echo "    Output file (relative input): ${SIGNATURE_FILE_REL}"
echo "    Signing file (absolute path): ${CHECKSUMS_FILE}"
echo "    Output file (absolute path): ${SIGNATURE_FILE}"


# --- Use a temporary file for the GPG passphrase ---
# Create a temporary file and store its name in PASSPHRASE_FILE
PASSPHRASE_FILE=$(mktemp)
echo "    Temporary passphrase file created at: ${PASSPHRASE_FILE}"

# Set up a trap to clean up the temporary passphrase file on exit of *this script*
trap "echo 'Cleaning up temporary passphrase file...' && rm -f \"${PASSPHRASE_FILE}\"" EXIT

# Write the passphrase from stdin (where it's piped to the script) to the temporary file
# Use 'cat -' to read from stdin
echo "    Reading passphrase from stdin..."
cat - > "${PASSPHRASE_FILE}"
echo "    Passphrase written to temporary file."

# --- Set GPG Environment Variables if needed ---
# GPG_TTY is often needed for non-interactive GPG operations
export GPG_TTY="$(tty)"
echo "    Set GPG_TTY=${GPG_TTY}"

# Export custom GNUPGHOME if provided
if [ -n "$GPG_HOME_PATH" ]; then
  export GNUPGHOME="$GPG_HOME_PATH"
  echo "    Set GNUPGHOME=${GNUPGHOME}"
fi

# --- Debug: Print GPG command details (using absolute paths) ---
echo "    DEBUG (script): Final GNUPGHOME value: '${GNUPGHOME:-<not set>}'"
echo "    DEBUG (script): Checksums file path (absolute): '${CHECKSUMS_FILE}'"
echo "    DEBUG (script): Signature output path (absolute): '${SIGNATURE_FILE}'"
echo "    DEBUG (script): Passphrase file: '${PASSPHRASE_FILE}'" # This is already an absolute path from mktemp

# --- Run GPG Command (Binary Signature) ---
echo "    Executing GPG command..."
# Print the exact command GPG will run (using absolute paths for input/output)
#echo "    DEBUG (script): gpg --batch --yes --passphrase-file \"${PASSPHRASE_FILE}\" --pinentry-mode loopback ${GPG_KEY_OPT} --output \"${SIGNATURE_FILE}\" --detach-sign \"${CHECKSUMS_FILE}\""

gpg --batch --yes --passphrase-file "${PASSPHRASE_FILE}" --pinentry-mode loopback ${GPG_KEY_OPT} --output "${SIGNATURE_FILE}" --detach-sign "${CHECKSUMS_FILE}"

# Check the exit status of the GPG command
GPG_EXIT_CODE=$?
if [ "$GPG_EXIT_CODE" -ne 0 ]; then
    echo "Error: GPG signing failed with exit code ${GPG_EXIT_CODE}."
    # The trap will handle temp file cleanup
    exit ${GPG_EXIT_CODE}
fi

echo "    GPG signing command completed successfully."

# The trap will handle temp file cleanup

echo "Signing script finished."
