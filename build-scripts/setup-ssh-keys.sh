#!/usr/bin/env bash

SCRIPTNAME=$(basename $0)
VERSION=1.0.1

# SCRIPT START
exec 1>&2

cat <<EOF
##########################
# $SCRIPTNAME $VERSION 
##########################
EOF

if [ "$CI" != "true" ] && [ "$CIRCLECI" != "true" ] && [ "$CI" != "TRUE" ] && [ "$CIRCLECI" != "TRUE" ]; then
  echo "INFO: We are NOT running in CI. Therefore there is no need to configure a global SSH key."
  echo ""
  exit 0
fi

echo "INFO: configuring global ssh key"
if [ -z "$DEVOPS_GLOBAL_SSH_PRIVATE_KEY_B64" ]; then
  echo "ERROR: Failed to load DevOps global SSH key: No DEVOPS_GLOBAL_SSH_PRIVATE_KEY_B64 environment variable is defined"
  echo "Please contact the DevOps team"
  echo ""
  exit 1
fi

GLOBAL_KEY_FILE=/tmp/.global_key_$RANDOM
echo ""
echo "INFO: Deleting all SSH keys loaded by the agent"
ssh-add -D

echo ""
echo "Importing DevOps global key from ${GLOBAL_KEY_FILE}"
echo "$DEVOPS_GLOBAL_SSH_PRIVATE_KEY_B64" | base64 -d >$GLOBAL_KEY_FILE
chmod 600 $GLOBAL_KEY_FILE
ssh-add $GLOBAL_KEY_FILE
RETVAL=$?
rm -f ${GLOBAL_KEY_FILE}

if [ ${RETVAL} -ne 0 ]; then
  echo "WARNING: ssh-add command returned error $RETVAL when trying to add the ssh key!"
  echo "In case cloning other private repositories fails as part of this build - please contact the DevOps team"
else
  echo "SUCCESS! You can now clone our private repositories"
fi

BITBUCKET_SSH_HOST_KEY_WORKAROUND_HEADER="# ${SCRIPTNAME}-bitbucket.org-ssh-host-key-workaround"
if [ -f ~/.ssh/config ] && grep -q "^#${BITBUCKET_SSH_HOST_KEY_WORKAROUND_HEADER}$" ~/.ssh/config; then
  echo "~/.ssh/config is already configured for bitbucket ssh host key"
else
  echo "Configuring ssh host key setting for bitbucket.org in ~/.ssh/config"
  test -d ~/.ssh || mkdir ~/.ssh
  echo -e "${BITBUCKET_SSH_HOST_KEY_WORKAROUND_HEADER}\nHost bitbucket.org\n  HostName bitbucket.org\n  CheckHostIP=no" >>~/.ssh/config
fi

echo ""
exit ${RETVAL}
