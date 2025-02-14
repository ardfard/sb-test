#!/usr/bin/env bash

API_URL="http://localhost:8080/upload"
TEST_AUDIO_FILE="test.m4a"

if [[ ! -f "$TEST_AUDIO_FILE" ]]; then
    echo "Test audio file '$TEST_AUDIO_FILE' not found. Please make sure it exists."
    exit 1
fi

echo "Uploading $TEST_AUDIO_FILE to $API_URL..."

# Use curl to send a POST request with the file as form-data
response=$(curl -s -w "\nHTTP_STATUS:%{http_code}\n" -F "audio=@${TEST_AUDIO_FILE}" "$API_URL")

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
