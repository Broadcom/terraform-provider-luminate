#!/bin/bash

VERSION=$(cat VERSION)

GITHUB_BASE_URL="https://api.github.com/repos/royeectu/${CIRCLE_PROJECT_REPONAME}"

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
    curl -f -X POST \
        -H "Content-Type: application/octet-stream" \
        -H "Authorization: token ${AUTH}" \
        --data-binary @"release/${FILE}" \
        "${UPLOAD_URL}?name=${FILE}"
        RETVAL=$?
        echo ""
        if [ $RETVAL -ne 0]
        then
          echo "Error! Failed to upload $FILE to $UPLOAD_URL - curl returned error #$RETVAL"
          exit $RETVAL
        fi
done

exit 0