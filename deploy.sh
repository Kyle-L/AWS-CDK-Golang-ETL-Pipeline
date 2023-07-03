#!/bin/bash
# This script deploys the CDK stack and frontend to AWS
#
# This script assumes you have the following installed:
# - AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
# - jq: https://stedolan.github.io/jq/download/
# - npm: https://www.npmjs.com/get-npm
# - cdk: https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html
#
# Copyright 2023 Kyle Lierer. All Rights Reserved.
#

# Gets the stack name from the CLI
STACK_NAME=$1

# Defines a few colors for use in the script
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# If no stack name is provided, exit
if [ -z "$STACK_NAME" ]; then
    echo -e "${RED}No stack name provided!${NC}"
    echo -e "${RED}Usage: ./deploy.sh <stack-name>${NC}"
    exit 1
fi

# Check if AWS credentials are configured
aws sts get-caller-identity > /dev/null
if [ $? -ne 0 ]; then
    echo -e "${RED}AWS credentials are not configured!${NC}"
    exit 1
fi

# Deploy the CDK stack to the AWS account/region
echo -e "${BLUE}Deploying the CDK workshop stack...${NC}"
cd backend
cdk deploy -c name=$STACK_NAME --outputs-file ./outputs.json

if [ $? -ne 0 ]; then
    echo -e "${RED}CDK stack deployment failed!${NC}"
    exit 1
fi

echo -e "${GREEN}CDK stack deployed!${NC}"

# Get the API Gateway URL from the CDK json output
API_URL=$(cat outputs.json | jq -r '."'$STACK_NAME'".ApiUrl')
BUCKET_NAME=$(cat outputs.json | jq -r '."'$STACK_NAME'".BucketName')

# Build the frontend with the API Gateway URL
echo -e "${BLUE}Building the frontend...${NC}"
cd ../frontend
API_ENDPOINT=$API_URL npm run build 
echo -e "${GREEN}Frontend built!${NC}"

# Deploy the frontend to S3 and make it public
echo -e "${BLUE}Deploying the frontend to S3...${NC}"
aws s3 sync ./out/ s3://$BUCKET_NAME --delete > /dev/null

# Check if s3 sync was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Frontend deployed!${NC}"
else
    echo -e "${RED}Frontend deployment failed!${NC}"
    exit 1
fi

echo -e "${GREEN}Deployment complete!${NC}"
echo ""
echo -e "${BLUE}API URL: $API_URL${NC}"
echo -e "${BLUE}Frontend URL: https://$BUCKET_NAME-frontend.s3.amazonaws.com${NC}"
