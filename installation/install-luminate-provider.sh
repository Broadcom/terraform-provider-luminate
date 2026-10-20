#!/usr/bin/env bash
# Copyright (c) Broadcom Inc.
# SPDX-License-Identifier: MPL-2.0

# Installs the luminate Terraform provider into the local Terraform plugins directory.
#
# Usage:
#   ./install-luminate-provider.sh [version]          e.g. ./install-luminate-provider.sh 1.2.18
#
# Without downloading the script first:
#   curl -fsSL https://raw.githubusercontent.com/Broadcom/terraform-provider-luminate/master/installation/install-luminate-provider.sh | bash
#   curl -fsSL <same url> | bash -s -- 1.2.18
#
# If no version is given, the latest release is installed.
# On Windows, run it from Git Bash (or WSL for the Linux build).
set -euo pipefail

REPO_URL="https://github.com/Broadcom/terraform-provider-luminate"

if [ $# -gt 1 ]; then
  echo "Usage: $0 [version]   (e.g. $0 1.2.18; defaults to the latest release)" >&2
  exit 1
fi

LUMINATE_TF_VERSION="${1:-}"

if [ -z "$LUMINATE_TF_VERSION" ]; then
  # /releases/latest redirects to /releases/tag/<version>
  LATEST_URL="$(curl -fsSL -o /dev/null -w '%{url_effective}' "${REPO_URL}/releases/latest")"
  LUMINATE_TF_VERSION="${LATEST_URL##*/}"
  if [ -z "$LUMINATE_TF_VERSION" ] || [ "$LUMINATE_TF_VERSION" = "latest" ]; then
    echo "Failed to determine the latest release version" >&2
    exit 1
  fi
  echo "No version specified, using latest release: ${LUMINATE_TF_VERSION}"
fi

# Detect OS and architecture
case "$(uname -s)" in
  Darwin) OS=darwin ;;
  Linux)  OS=linux ;;
  MINGW*|MSYS*|CYGWIN*) OS=windows ;;
  *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  arm64|aarch64) ARCH=arm64 ;;
  x86_64|amd64)  ARCH=amd64 ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

PLATFORM="${OS}_${ARCH}"
if [ "$PLATFORM" = "linux_arm64" ]; then
  echo "No linux_arm64 build is published for the luminate provider" >&2
  exit 1
fi

EXT=""
PLUGINS_DIR="$HOME/.terraform.d/plugins"
if [ "$OS" = "windows" ]; then
  EXT=".exe"
  # Terraform on Windows reads plugins from %APPDATA%\terraform.d\plugins
  PLUGINS_DIR="$(cygpath -u "$APPDATA")/terraform.d/plugins"
fi

ZIP="terraform-provider-luminate_${LUMINATE_TF_VERSION}_${PLATFORM}.zip"
SUMS="terraform-provider-luminate_${LUMINATE_TF_VERSION}_SHA256SUMS"
BASE_URL="${REPO_URL}/releases/download/${LUMINATE_TF_VERSION}"
DEST="${PLUGINS_DIR}/broadcom.com/broadcom/luminate/${LUMINATE_TF_VERSION}/${PLATFORM}"

# sha256sum on Linux, shasum on macOS
if command -v sha256sum >/dev/null 2>&1; then
  sha256() { sha256sum "$1" | awk '{print $1}'; }
else
  sha256() { shasum -a 256 "$1" | awk '{print $1}'; }
fi

extract_zip() {
  if command -v unzip >/dev/null 2>&1; then
    unzip -q -o "$1" -d "$2"
  else
    # Git Bash does not ship unzip
    powershell.exe -NoProfile -Command "Expand-Archive -Force -LiteralPath '$(cygpath -w "$1")' -DestinationPath '$(cygpath -w "$2")'"
  fi
}

echo "Installing luminate provider ${LUMINATE_TF_VERSION} (${PLATFORM}) into ${DEST}"

# Work in a temp dir so an interrupted install never leaves a partial provider in the plugins dir
TMP_DIR="$(mktemp -d)"
trap 'rm -rf -- "$TMP_DIR"' EXIT

curl -fsSL -o "${TMP_DIR}/${ZIP}" "${BASE_URL}/${ZIP}"
curl -fsSL -o "${TMP_DIR}/${SUMS}" "${BASE_URL}/${SUMS}"

# Entries in SHA256SUMS are prefixed with a path (dist/<version>/), so match on the file name only
EXPECTED="$(awk -v f="$ZIP" '{n=split($2, p, "/"); if (p[n] == f) print $1}' "${TMP_DIR}/${SUMS}")"
if [ -z "$EXPECTED" ]; then
  echo "No checksum found for ${ZIP} in ${SUMS}" >&2
  exit 1
fi
ACTUAL="$(sha256 "${TMP_DIR}/${ZIP}")"
if [ "$EXPECTED" != "$ACTUAL" ]; then
  echo "Checksum mismatch for ${ZIP}: expected ${EXPECTED}, got ${ACTUAL}" >&2
  exit 1
fi
echo "Checksum verified"

extract_zip "${TMP_DIR}/${ZIP}" "$TMP_DIR"
# The release zip ships the binary with the version suffix; Terraform expects the plain name
BIN="terraform-provider-luminate${EXT}"
mv -f "${TMP_DIR}/terraform-provider-luminate_${LUMINATE_TF_VERSION}${EXT}" "${TMP_DIR}/${BIN}"
chmod +x "${TMP_DIR}/${BIN}"

# Remove the quarantine attribute (macOS only; ignore if not set)
if [ "$OS" = "darwin" ]; then
  xattr -d com.apple.quarantine "${TMP_DIR}/${BIN}" 2>/dev/null || true
fi

mkdir -p "$DEST"
mv -f "${TMP_DIR}/${BIN}" "${DEST}/${BIN}"

echo "Done: ${DEST}/${BIN}"
