#!/bin/bash
PACKAGE_NAME="shrampybot-dev"
PACKAGE_PATH="packages/shrampybot-dev/api"
SCRIPT_PATH="${PWD}"
REGION="eu-central-1"
BUCKET_NAME=$PACKAGE_NAME
echo -n "Making s3 bucket ${BUCKET_NAME}..."
aws s3 mb s3://shrampybot-dev 1>/dev/null 2>/dev/null
if [[ $? -eq 0 ]]; then
    echo "created!"
elif [[ $? -eq 1 ]]; then
    echo "already created."
fi

echo "Building ${PACKAGE_NAME}..."
cd "${PACKAGE_PATH}"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
GO_BUILD_RESULT=$?
cd "${SCRIPT_PATH}"
if [[ $GO_BUILD_RESULT -ne 0 ]]; then
    echo "Go build failed."
    exit $GO_BUILD_RESULT
fi

echo "Zipping function directory..."
rm -f "${SCRIPT_PATH}/${PACKAGE_NAME}.zip"
cd "${PACKAGE_PATH}"
zip -r "${SCRIPT_PATH}/${PACKAGE_NAME}.zip" ./* 1>/dev/null
cd "${SCRIPT_PATH}"

echo "Deploying ${PACKAGE_NAME}..."
aws lambda update-function-code --zip-file "fileb://${PACKAGE_NAME}.zip" \
    --function-name "${PACKAGE_NAME}" --publish --no-paginate \
    --no-cli-pager --no-cli-auto-prompt 1>/dev/null

