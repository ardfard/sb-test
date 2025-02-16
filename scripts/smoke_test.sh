#!/usr/bin/env bash

# Smoke Test Script for Audio Upload and Download API
# This script tests:
# 1. /audio/user/{user_id}/phrase/{phrase_id} (Upload)
# 2. /audio/user/{user_id}/phrase/{phrase_id}/{format} (Download)
#
# Prerequisites:
# - Ensure your API server is running on http://localhost:8080
# - Place a test audio file named "test.m4a" in the same directory as this script
#
# Usage:
#   chmod +x smoke_test.sh
#   ./smoke_test.sh

API_HOST="http://localhost:8080"
TEST_AUDIO_FILE="tests/fixtures/test.m4a"
DOWNLOAD_DIR="downloaded"

# Create download directory if it doesn't exist
mkdir -p "$DOWNLOAD_DIR"

if [[ ! -f "$TEST_AUDIO_FILE" ]]; then
    echo "Test audio file '$TEST_AUDIO_FILE' not found. Please make sure it exists."
    exit 1
fi

# Create user
USER_ID=$(curl -s -X POST -H "Content-Type: application/json" \
    -d '{"name": "Ashen One"}' \
    "${API_HOST}/users" | jq -r '.id')

# Create phrase
PHRASE_ID=$(curl -s -X POST -H "Content-Type: application/json" \
    -d '{"text": "Only in truth, the Lords will abandon their thrones, and the Unkindled will rise"}' \
    "${API_HOST}/users/${USER_ID}/phrases" | jq -r '.id')

# Test Upload
echo "Testing Upload Endpoint..."
UPLOAD_ENDPOINT="${API_HOST}/audio/user/${USER_ID}/phrase/${PHRASE_ID}"

echo "Uploading $TEST_AUDIO_FILE to $UPLOAD_ENDPOINT..."

upload_response=$(curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
    -F "audio_file=@${TEST_AUDIO_FILE}" \
    "${UPLOAD_ENDPOINT}")

upload_status=$(echo "$upload_response" | sed -n 's/HTTP_STATUS://p' | tr -d ' ')
upload_body=$(echo "$upload_response" | sed '/HTTP_STATUS:/d')

echo "Upload Response HTTP Status: $upload_status"
echo "Upload Response Body:"
echo "$upload_body"

if [ "$upload_status" -ne 200 ]; then
    echo "Upload test failed with status code $upload_status"
    exit 1
fi

# Wait for conversion (adjust time as needed)
echo "Waiting for conversion to complete..."
sleep 2

# Test Download
echo "Testing Download Endpoint..."
DOWNLOAD_ENDPOINT="${API_HOST}/audio/user/${USER_ID}/phrase/${PHRASE_ID}/m4a"
DOWNLOADED_FILE="${DOWNLOAD_DIR}/audio_${USER_ID}_${PHRASE_ID}.m4a"

echo "Downloading converted audio from $DOWNLOAD_ENDPOINT..."

download_response=$(curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
    -o "$DOWNLOADED_FILE" \
    "${DOWNLOAD_ENDPOINT}")

download_status=$(echo "$download_response" | sed -n 's/HTTP_STATUS://p' | tr -d ' ')

echo "Download Response HTTP Status: $download_status"

if [ "$download_status" -ne 200 ]; then
    echo "Download test failed with status code $download_status"
    exit 1
fi

if [ ! -f "$DOWNLOADED_FILE" ]; then
    echo "Downloaded file not found"
    exit 1
fi

# Check if downloaded file is a valid m4a file
file_type=$(file -b "$DOWNLOADED_FILE")
if [[ ! "$file_type" =~ "M4A" ]]; then
    echo "Downloaded file is not a valid M4A file. File type: $file_type"
    exit 1
fi

echo "Downloaded file saved to: $DOWNLOADED_FILE"
echo "File type: $file_type"
echo "File size: $(ls -lh "$DOWNLOADED_FILE" | awk '{print $5}')"

echo "Smoke test succeeded!" 
