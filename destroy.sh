#!/bin/bash
# This script destroy the CDK stack and frontend to AWS
#
# This script assumes you have the following installed:
# - AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
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

# Deploy the CDK stack to the AWS account/region
echo -e "${BLUE}Destroying the CDK workshop stack...${NC}"
cd backend
cdk destroy -c name=$STACK_NAME
echo -e "${GREEN}CDK stack destroy!${NC}"