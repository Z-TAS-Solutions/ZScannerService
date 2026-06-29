#!/bin/bash

# Configuration
REPO="Z-TAS-Solutions/ZScannerService"
BINARY_NAME="ZScannerService"
TARGET_DIR="/home/ZTAS"
TOKEN=$(cat "$TARGET_DIR/.github_token")

# Ensuring da directory exists..
mkdir -p "$TARGET_DIR"

# Wait for a basic network connection so curl doesn't fail instantly
echo "Checking network connectivity..."
until curl -s --connect-timeout 3 https://api.github.com > /dev/null; do
    echo "Waiting for network..."
    sleep 2
done

# Fetching latest successful workflow run details
RUN_DATA=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "https://api.github.com/repos/$REPO/actions/runs?branch=main&status=success&per_page=1")

# Extract the Run ID and the Commit Hash
LATEST_RUN_ID=$(echo "$RUN_DATA" | jq -r '.workflow_runs[0].id')
LATEST_COMMIT=$(echo "$RUN_DATA" | jq -r '.workflow_runs[0].head_sha')

if [ "$LATEST_RUN_ID" == "null" ] || [ -z "$LATEST_RUN_ID" ]; then
    echo "No successful workflow builds found."
    exit 0
fi

# Read local version if the thing exists that is
CURRENT_COMMIT=""
[ -f "$TARGET_DIR/current_commit.txt" ] && CURRENT_COMMIT=$(cat "$TARGET_DIR/current_commit.txt")

# If there is a new commit build, fetch the artifact download URL
if [ "$LATEST_COMMIT" != "$CURRENT_COMMIT" ]; then
    echo "New build found for commit: ${LATEST_COMMIT:0:7}"

    ARTIFACT_DATA=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "https://api.github.com/repos/$REPO/actions/runs/$LATEST_RUN_ID/artifacts")
    
    DOWNLOAD_URL=$(echo "$ARTIFACT_DATA" | jq -r '.artifacts[] | select(.name=="ZScannerService-artifact") | .archive_download_url')

    if [ -n "$DOWNLOAD_URL" ] && [ "$DOWNLOAD_URL" != "null" ]; then
        curl -L -H "Authorization: Bearer $TOKEN" -o "$TARGET_DIR/artifact.zip" "$DOWNLOAD_URL"
        
        unzip -o "$TARGET_DIR/artifact.zip" -d "$TARGET_DIR/"
        rm -f "$TARGET_DIR/artifact.zip"

        chmod +x "$TARGET_DIR/$BINARY_NAME"
        echo "$LATEST_COMMIT" > "$TARGET_DIR/current_commit.txt"
        echo "Update applied successfully."
    else
        echo "Could not find matching artifact name in this run."
    fi
else
    echo "Already running the latest build (${CURRENT_COMMIT:0:7})."
fi
