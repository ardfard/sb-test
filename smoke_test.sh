#!/usr/bin/env bash

# Smoke Test Script for Audio Upload API
# This script tests the /audio/user/{user_id}/phrase/{phrase_id} endpoint
#
# Prerequisites:
# - Ensure your API server is running on http://localhost:8080
# - Place a test audio file named "test.m4a" in the same directory as this script
#
# Usage:
#   chmod +x smoke_test.sh
#   ./smoke_test.sh

API_HOST="http://localhost:8080"
TEST_AUDIO_FILE="test.m4a"
USER_ID="123"
PHRASE_ID="456"

if [[ ! -f "$TEST_AUDIO_FILE" ]]; then
    echo "Test audio file '$TEST_AUDIO_FILE' not found. Please make sure it exists."
    exit 1
fi

ENDPOINT="${API_HOST}/audio/user/${USER_ID}/phrase/${PHRASE_ID}"

echo "Uploading $TEST_AUDIO_FILE to $ENDPOINT..."

# Use curl to send a POST request with the file as form-data
response=$(curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
    -F "audio=@${TEST_AUDIO_FILE}" \
    "${ENDPOINT}")

# Extract HTTP status code and body from the response
http_status=$(echo "$response" | sed -n 's/HTTP_STATUS://p' | tr -d ' ')
body=$(echo "$response" | sed '/HTTP_STATUS:/d')

echo "Response HTTP Status: $http_status"
echo "Response Body:"
echo "$body"

if [ "$http_status" -ne 200 ]; then
    echo "Smoke test failed with status code $http_status."
    exit 1
fi

echo "Smoke test succeeded." 
