#!/bin/bash

# This script configures SSH keys from environment variables to clone a private repository.
# It's designed to be robust for use in CI/CD pipelines.

# Exit immediately if a command exits with a non-zero status.
set -e

col_red=$esc_seq"31;01m"

echo "--> Setting up SSH keys to clone API documentation..."

# 1. Validate that the required environment variables are set
	echo "INFO: Validating SSH key environment variables before writing to files..."

	if [[ -z "${BROADCOM_GITHUB_ACCESS_LUMINATE_PRIVATE_KEY_B64}" ]]; then
		echo -e "${col_red}ERROR: BROADCOM_GITHUB_ACCESS_LUMINATE_PRIVATE_KEY_B64 is not set.${col_reset}"
		exit 1
	fi

	if [[ -z "${BROADCOM_GITHUB_ACCESS_GITHUB_PRIVATE_KEY_B64}" ]]; then
		echo -e "${col_red}ERROR: BROADCOM_GITHUB_ACCESS_GITHUB_PRIVATE_KEY_B64 is not set.${col_reset}"
		exit 1
	fi

# 2. Create the .ssh directory and set the correct, secure permissions
echo "--> Configuring ~/.ssh directory..."
mkdir -p ~/.ssh
chmod 700 ~/.ssh

echo "--> Writing decoded SSH keys..."
echo "$BROADCOM_GITHUB_ACCESS_LUMINATE_PRIVATE_KEY_B64" | base64 -d > ~/.ssh/id_rsa
echo "$BROADCOM_GITHUB_ACCESS_GITHUB_PRIVATE_KEY_B64" | base64 -d > ~/.ssh/id_ed25519

# 4. Set the correct, restrictive permissions on the key files, as required by SSH.
chmod 600 ~/.ssh/id_rsa ~/.ssh/id_ed25519

#    prompts that would cause a CI build to hang or fail.
echo "--> Configuring SSH for github.gwd.broadcom.net..."
ssh-keyscan -t rsa broadcom-github.ssh.luminate.luminatesite.com >> ~/.ssh/known_hosts
cat <<EOF > ~/.ssh/config
Host github.gwd.broadcom.net
  HostName broadcom-github.ssh.luminate.luminatesite.com
  User git@broadcom-github
  IdentityFile ~/.ssh/id_ed25519
  AddKeysToAgent yes
  ForwardAgent yes
  IdentitiesOnly yes
EOF

# Configure Git to use SSH for github.gwd.broadcom.net
echo "--> Configuring Git to use SSH for github.gwd.broadcom.net..."
git config --global url."git@github.gwd.broadcom.net:".insteadOf "https://github.gwd.broadcom.net/"

# 5. Start the ssh-agent in the background and add the keys to it.
echo "--> Starting ssh-agent and adding keys..."
eval "$(ssh-agent -s)"
trap "echo 'Killing ssh-agent (PID: $SSH_AGENT_PID)...' && kill $SSH_AGENT_PID" EXIT
ssh-add ~/.ssh/id_rsa
ssh-add ~/.ssh/id_ed25519
ssh-add -l

# 6. Check connection to github via SAC
echo "INFO: Checking connection to github.gwd.broadcom.net"
ssh -T github.gwd.broadcom.net

# 7. Clone the repository using the specified key.
echo "--> Cloning ztna-api-documentation repository..."
git clone github.gwd.broadcom.net:SED/ztna-api-documentation.git

LUM_API_DOC_REPO=ztna-api-documentation
cd ${LUM_API_DOC_REPO} || return 1
echo "INFO: creating go.mod for module ${LUM_API_DOC_REPO}"
go mod init ${LUM_API_DOC_REPO} || return $?
cd go/sdk
go mod init ${LUM_API_DOC_REPO} || return $?
cd /home/circleci/project

echo "--> Replacing Go module to point to the local clone..."
go mod edit -replace github.gwd.broadcom.net/SED/${LUM_API_DOC_REPO}=./${LUM_API_DOC_REPO} || return $?
go mod edit -replace github.gwd.broadcom.net/SED/${LUM_API_DOC_REPO}/go/sdk=./${LUM_API_DOC_REPO}/go/sdk || return $?
go mod tidy

echo "--> Successfully configured Go module replacement."
