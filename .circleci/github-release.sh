#!/bin/bash

VERSION=$(cat VERSION)

GITHUB_BASE_URL="https://api.github.com/repos/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}"

AUTH="${GITHUB_ACCESS_TOKEN}"


read -r -d '' PAYLOAD << EOF
{
    "tag_name": "${VERSION}",
    "name": "${VERSION}"
}
EOF

RESP=$(curl -X POST -H "Authorization: token ${AUTH}" ${GITHUB_BASE_URL}/releases -d "${PAYLOAD}")

UPLOAD_URL=$(echo ${RESP} | jq -r .upload_url | cut -f1 -d{ )

for FILE in $(ls release)
do
    echo "Uploading $FILE"
    curl -X POST \
        -H "Content-Type: application/octet-stream" \
        --data-binary @"release/${FILE}" \
        "${UPLOAD_URL}?name=${FILE}&${AUTH}"
done
