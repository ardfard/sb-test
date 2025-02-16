#!/usr/bin/env bash

set -euo pipefail

# Script configuration
IMAGE_NAME="sb-test"
IMAGE_TAG="latest"
CONTAINER_NAME="sb-test-instance"
CONFIG_PATH="config.yaml"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# make sure the current directory is the root of the project
cd "$(dirname "$0")/.."

echo -e "${GREEN}Building Docker image...${NC}"
docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" .

echo -e "${GREEN}Starting container...${NC}"
docker run -d \
    --rm \
    --name "${CONTAINER_NAME}" \
    -p 8080:8080 \
    "${IMAGE_NAME}:${IMAGE_TAG}"

# Wait for the service to be ready
echo -e "${GREEN}Waiting for service to be ready...${NC}"
sleep 5

echo -e "${GREEN}Running smoke test...${NC}"
./scripts/smoke_test.sh

# Cleanup
echo -e "${GREEN}Cleaning up...${NC}"
docker stop "${CONTAINER_NAME}"
docker rm "${CONTAINER_NAME}"

echo -e "${GREEN}Done!${NC}" 