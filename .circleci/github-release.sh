#!/bin/bash

VERSION=$(cat VERSION)

GITHUB_BASE_URL="https://api.github.com/repos/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}"

AUTH="${GITHUB_ACCESS_TOKEN}"

read -r -d '' PAYLOAD <<EOF
{
    "tag_name": "${VERSION}",
    "name": "${VERSION}"
}
EOF

echo "Creating release $VERSION on URL: $GITHUB_BASE_URL"
OUTPUT=$(mktemp)
RESP=$(curl -f -X POST -H "Authorization: token ${AUTH}" ${GITHUB_BASE_URL}/releases -d "${PAYLOAD}" -w "%{response_code}" -o $OUTPUT)
if [ $? -ne 0 ]; then
  echo "ERROR: creating the release failed with HTTP/$RESP"
  cat "$OUTPUT"
  exit 1
fi

UPLOAD_URL=$(cat ${OUTPUT} | jq -r .upload_url | cut -f1 -d{)
if [ -z "$UPLOAD_URL" ] || [ "$UPLOAD_URL" == "null" ]; then
  echo "ERROR: failed to parse upload URL as json (resolved as \"${UPLOAD_URL}\") from body:"
  cat ${OUTPUT}
  exit 1
fi

echo "Upload URL: $UPLOAD_URL"

for FILE in $(ls release); do
  echo "Uploading $FILE"
  curl -f -X POST \
    -H "Content-Type: application/octet-stream" \
    -H "Authorization: token ${AUTH}" \
    --data-binary @"release/${FILE}" \
    "${UPLOAD_URL}?name=${FILE}"
  RETVAL=$?
  echo ""
  if [ $RETVAL -ne 0 ]; then
    echo "Error! Failed to upload $FILE to $UPLOAD_URL - curl returned error #$RETVAL"
    exit $RETVAL
  fi
done

exit 0
