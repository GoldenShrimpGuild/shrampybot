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
    echo "already created or error occurred."
fi

echo "Building ${PACKAGE_NAME}..."
cd "${PACKAGE_PATH}"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
cd "${SCRIPT_PATH}"

echo "Zipping function directory..."
rm -f "${SCRIPT_PATH}/${PACKAGE_NAME}.zip"
cd "${PACKAGE_PATH}"
zip -r "${SCRIPT_PATH}/${PACKAGE_NAME}.zip" ./*
cd "${SCRIPT_PATH}"

echo "Deploying ${PACKAGE_NAME}..."
aws lambda update-function-code --zip-file "fileb://${PACKAGE_NAME}.zip" \
    --function-name "${PACKAGE_NAME}" --publish

# aws cloudformation package --template-file template-dev.yml --s3-bucket $BUCKET_NAME --output-template-file out-dev.yml 1>/dev/null
# aws cloudformation deploy --template-file out-dev.yml --stack-name shrampybot-dev --capabilities CAPABILITY_NAMED_IAM 1>/dev/null
