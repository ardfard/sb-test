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

# loop until the service is ready
while ! curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health; do
    sleep 1
done

echo -e "${GREEN}Service is running at http://localhost:8080...${NC}"

echo -e "${GREEN}Running smoke test...${NC}"
./scripts/smoke_test.sh

# cleanup the container if its running in CI environment
if [ -n "${CI:-}" ]; then
    echo -e "${GREEN}Cleaning up container...${NC}"
    docker rm -f "${CONTAINER_NAME}"
fi

echo -e "${GREEN}Smoke test completed successfully!${NC}"
echo -e "${GREEN}You can now test the service using the following command:${NC}"
echo -e "${GREEN}curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{\"name\": \"John Doe\"}'${NC}"
echo -e "${GREEN}curl -X POST http://localhost:8080/users/{user_id}/phrases -H 'Content-Type: application/json' -d '{\"text\": \"Hello, world!\"}'${NC}"
echo -e "${GREEN}curl -X POST http://localhost:8080/audio/user/{user_id}/phrase/{phrase_id} -H 'Content-Type: multipart/form-data' -F 'file=@path/to/your/audio/file'${NC}"
echo -e "${GREEN}curl -X GET http://localhost:8080/audio/user/{user_id}/phrase/{phrase_id}/{format}${NC}"
