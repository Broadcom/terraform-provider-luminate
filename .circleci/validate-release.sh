#!/bin/bash

VERSION=$(cat VERSION)

GITHUB_BASE_URL="https://api.github.com/repos/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}"

AUTH="access_token=${GITHUB_ACCESS_TOKEN}"

curl "${GITHUB_BASE_URL}/releases?${AUTH}" | jq -r .[].name | grep ${VERSION}
if [[ $? == 0 ]]
then
  echo "Error: Release ${VERSION} already exists. Did you forget to increment version number?"
  exit 1
fi
