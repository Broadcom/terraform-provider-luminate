#!/bin/bash

if [ -z "$CIRCLE_PROJECT_USERNAME" ]
then
  echo "Error: The environment variable CIRCLE_PROJECT_USERNAME is not defined! Cannot determine Github base URL!"
  exit 1
fi

if [ -z "CIRCLE_PROJECT_REPONAME" ]
then
  echo "Error: The environment variable CIRCLE_PROJECT_REPONAME is not defined! Cannot determine Github base URL!"
  exit 2
fi

VERSION=$(cat VERSION)
GITHUB_BASE_URL="https://api.github.com/repos/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}"
TMPFILE=/tmp/`basename ${0}`.tmp
test -f "$TMPFILE" && rm -f "$TMPFILE"

echo "Getting releases from ${GITHUB_BASE_URL}/releases using tmp file $TMPFILE"
curl -sSf -m10 "${GITHUB_BASE_URL}/releases" -o $TMPFILE
RETVAL=$?
if [ ${RETVAL} -ne 0 ]
then
  echo "Error: Failed to validate that the version $VERSION is unique - curl returned error #${RETVAL}"
  exit 3
fi

if [ ! -s "$TMPFILE" ]
then
  echo "Error: curl returned no output"
  exit 4
fi

RELEASES=$(cat "$TMPFILE" 2>/dev/null | jq -r .[].name)
if [ -z "$RELEASES" ]
then
  echo "Error: Failed to parse response from $TMPFILE"
  echo "TMPFILE:"
  cat "$TMPFILE"
  echo "JSON Parse:"
  cat "$TMPFILE" | jq -r .[].name
  exit 5
fi

echo "$RELEASES" | grep ${VERSION}
if [[ $? == 0 ]]
then
  echo "Error: Release ${VERSION} already exists. Did you forget to increment version number?"
  exit 5
else
  echo "Ok: Release ${VERSION} is unique"
  exit 0
fi
