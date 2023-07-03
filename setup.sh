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

# Defines a few colors for use in the script
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Checks if jq is installed
jq --version > /dev/null
if [ $? -ne 0 ]; then
    echo -e "${RED}jq is not installed. Please install it before continuing.${NC}"
    exit 1
else
    echo "${GREEN}jq is installed.${NC}"
fi

# Checks to see if golang 1.18 is installed
version=$(go version)
if [[ $version != *"1.18"* ]]; then
    echo -e "${RED}Go 1.18 is not installed. Please install it before continuing.${NC}"
    exit 1
else
    echo "${GREEN}Go 1.18 is installed.${NC}"
fi

# Checks if Node 18
version=$(node --version)
if [[ $version != *"v18"* ]]; then
    echo "${RED}Node 18 is not installed. Please install it before continuing.${NC}"
    exit 1
else
    echo "${GREEN}Node 18 is installed.${NC}"
fi

echo "Installing frontend dependencies..."
cd frontend
npm install
echo -e "${GREEN}Frontend dependencies installed!${NC}"

echo "Installing backend dependencies..."
cd ../backend
go mod download
echo -e "${GREEN}Backend dependencies installed!${NC}"

echo "Done!"