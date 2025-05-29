#!/bin/bash

# Get the provider version from the VERSION file
PROVIDER_NAME=terraform-provider-luminate
VERSION=$(cat VERSION)
if [ -z "$VERSION" ]; then
  echo "Error: VERSION file is empty or not found. Please ensure the VERSION file exists and contains a valid version string."
  exit 1
fi

RELEASE_DIR=dist/$VERSION
MANIFEST_FILENAME="${PROVIDER_NAME}_${VERSION}_manifest.json"
# Get the version string without the 'v' for filename matching/parsing
VERSION_NO_V=$(echo "$VERSION" | sed 's/^v//')

# --- Validation ---
# Directory validation is done after defining RELEASE_DIR
if [ ! -d "$RELEASE_DIR" ]; then
  echo "Error: Release directory not found: ${RELEASE_DIR}. Did the build complete successfully?"
  exit 1
fi

# Create an array to hold the platform entries
PLATFORMS=()

# Check if shasum or sha256sum is available
CHECKSUM_CMD=""
if command -v shasum > /dev/null 2>&1; then
  CHECKSUM_CMD="shasum -a 256"
elif command -v sha256sum > /dev/null 2>&1; then
  CHECKSUM_CMD="sha256sum"
else
  echo "Error: Neither shasum nor sha256sum found. Cannot calculate checksums."
  exit 1
fi

echo "--> Generating manifest file: ${MANIFEST_FILENAME}"
echo "    Using release directory: ${RELEASE_DIR}"
echo "    Provider: ${PROVIDER_NAME}, Version: ${VERSION}"

# --- Build Platforms Array (as shell array of JSON strings) ---
PLATFORMS=()

# Navigate to the release directory to work with relative paths for filenames
# Store the original directory to return later
ORIG_DIR="$(pwd)"
cd "$RELEASE_DIR"
# Find and process zip files
# Need to handle case where no zip files are found
zip_files=(*.zip)
if [ "${zip_files[0]}" = "*.zip" ]; then
    echo "Error: No zip files found in ${RELEASE_DIR}"
    exit 1
fi

echo "    Found zip files:"
printf "      - %s\n" "${zip_files[@]}" # Print found zip files

for zip_file in "${zip_files[@]}"; do

  if [ ! -f "$zip_file" ]; then continue; fi # Skip if the pattern didn't match any files

  echo "    Processing zip file: ${zip_file}"

  OS=""
  ARCH=""
  ZIP_SHASUM=""

  # Calculate checksum for this specific zip file
  ZIP_SHASUM=$(${CHECKSUM_CMD} "${zip_file}" | awk '{print $1}')
  echo "      Calculated SHA: ${ZIP_SHASUM}"

 case "$zip_file" in
    ${PROVIDER_NAME}_${VERSION_NO_V}_darwin_amd64.zip)
      OS="darwin"
      ARCH="amd64"
      ;;
    ${PROVIDER_NAME}_${VERSION_NO_V}_darwin_arm64.zip)
      OS="darwin"
      ARCH="arm64"
      ;;
    ${PROVIDER_NAME}_${VERSION_NO_V}_linux_amd64.zip)
      OS="linux"
      ARCH="amd64"
      ;;
    ${PROVIDER_NAME}_${VERSION_NO_V}_windows_amd64.zip)
      OS="windows"
      ARCH="amd64"
      ;;
    *)
      echo "      Warning: Unrecognized filename format for OS/ARCH: ${zip_file}. Skipping platform."
      continue # Skip to the next zip_file
      ;;
  esac
  echo "      Matched OS: '${OS}', ARCH: '${ARCH}'"

  PLATFORM_JSON=$(printf '{\n      \"os\": \"%s\",\n      \"arch\": \"%s\",\n      \"filename\": \"%s\",\n      \"shasum\": \"%s\"\n    }' \
                    "$OS" "$ARCH" "$zip_file" "$ZIP_SHASUM")

  # Add the JSON string to the shell array
  PLATFORMS+=("$PLATFORM_JSON")

done

# Return to original directory
# Use stored original directory path for safety
cd "$ORIG_DIR"

JOINED_CONTENT=""
if [ ${#PLATFORMS[@]} -gt 0 ]; then
        first=true
        for item in "${PLATFORMS[@]}"; do
            INDENTED_PLATFORM_OBJECT=$(echo "$item" | sed 's/^/    /')
            if [ "$first" = true ]; then
              JOINED_CONTENT+="$INDENTED_PLATFORM_OBJECT"
              first=false
            else
              JOINED_CONTENT+=$(printf ",\n%s" "$INDENTED_PLATFORM_OBJECT")
            fi
        done
fi

PLATFORMS_JSON="  [
  ${JOINED_CONTENT}
    ]"

echo "    Finished processing zip files. PLATFORMS array has ${#PLATFORMS[@]} elements."

# --- Assemble Final JSON ---

echo "    Assembling final JSON..."

# Use printf to assemble the final JSON structure
# Output to the manifest file - Construct the full path here
MANIFEST_OUTPUT_PATH="${RELEASE_DIR}/${MANIFEST_FILENAME}" # Construct the full path here before printf

# Join the array elements with ",<\n>".
# Use printf to assemble the final JSON structure.
# This controls overall indentation and comma placement.
printf '{
  \"version\": \"%s\",
  \"protocols\":  [\"5.0\"],
  \"platforms\": [
    %s
  ]
}\n' \
          "$VERSION_NO_V" "$PLATFORMS_JSON" > "${MANIFEST_OUTPUT_PATH}"


echo "    Successfully created manifest file: ${MANIFEST_OUTPUT_PATH}"

echo "Manifest JSON creation complete."