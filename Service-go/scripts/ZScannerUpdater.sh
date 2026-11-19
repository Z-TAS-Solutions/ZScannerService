#!/bin/bash

# Configuration
REPO="Z-TAS-Solutions/ZScannerService"
TARGET_DIR="/home/zischl/ZTAS/ZScannerService"
TOKEN=$(cat "$TARGET_DIR/.github_token")

# Ensuring the target directory exists
mkdir -p "$TARGET_DIR"

# Wait for a basic network connection so curl doesn't fail instantly
echo "Checking network connectivity..."
until curl -s --connect-timeout 3 https://api.github.com > /dev/null; do
    echo "Waiting for network..."
    sleep 2
done

# Fetching latest successful workflow run details
RUN_DATA=$(curl -s -H "Authorization: Bearer $TOKEN" \
    "https://api.github.com/repos/$REPO/actions/runs?status=completed&per_page=1")

# Extract the Run ID and the Commit Hash
LATEST_RUN_ID=$(echo "$RUN_DATA" | jq -r '.workflow_runs[0].id')
LATEST_COMMIT=$(echo "$RUN_DATA" | jq -r '.workflow_runs[0].head_sha')

if [ "$LATEST_RUN_ID" == "null" ] || [ -z "$LATEST_RUN_ID" ]; then
    echo "No successful workflow builds found."
    exit 0
fi

# Read local version if it exists
CURRENT_COMMIT=""
[ -f "$TARGET_DIR/current_commit.txt" ] && CURRENT_COMMIT=$(cat "$TARGET_DIR/current_commit.txt")

# If there is a new commit build, fetch artifacts
if [ "$LATEST_COMMIT" != "$CURRENT_COMMIT" ]; then
    echo "New build found for commit: ${LATEST_COMMIT:0:7}"

    # Get artifact list for the latest run
    ARTIFACT_DATA=$(curl -s -H "Authorization: Bearer $TOKEN" \
        "https://api.github.com/repos/$REPO/actions/runs/$LATEST_RUN_ID/artifacts")
    
    UPDATED_ANY=false

    # Loop through both the Go and C++ artifact variants
    for VARIANT in "go" "cpp"; do
        ARTIFACT_NAME="ZScannerService-${VARIANT}-artifact"
        BINARY_NAME="zscanner-${VARIANT}"
        
        DOWNLOAD_URL=$(echo "$ARTIFACT_DATA" | jq -r --arg name "$ARTIFACT_NAME" '.artifacts[] | select(.name==$name) | .archive_download_url')

        if [ -n "$DOWNLOAD_URL" ] && [ "$DOWNLOAD_URL" != "null" ]; then
            echo "Downloading $ARTIFACT_NAME..."
            curl -L -H "Authorization: Bearer $TOKEN" -o "$TARGET_DIR/artifact.zip" "$DOWNLOAD_URL"
            
            # Extract to temporary directory to handle changing the file name cleanly
            mkdir -p "$TARGET_DIR/tmp"
            unzip -o "$TARGET_DIR/artifact.zip" -d "$TARGET_DIR/tmp/"
            rm -f "$TARGET_DIR/artifact.zip"

            if [ -f "$TARGET_DIR/tmp/ZScannerService" ]; then
                mv "$TARGET_DIR/tmp/ZScannerService" "$TARGET_DIR/$BINARY_NAME"
                chmod +x "$TARGET_DIR/$BINARY_NAME"
                echo "Successfully saved as: $BINARY_NAME"
                UPDATED_ANY=true
            else
                echo "Warning: Expected binary missing inside $ARTIFACT_NAME zip."
            fi
            rm -rf "$TARGET_DIR/tmp"
        else
            echo "Skipping $ARTIFACT_NAME (Not found in this pipeline run)."
        fi
    done

    # Save tracking commit only if at least one variant successfully extracted
    if [ "$UPDATED_ANY" = true ]; then
        echo "$LATEST_COMMIT" > "$TARGET_DIR/current_commit.txt"
        echo "Updates applied successfully."
    else
        echo "No valid binaries were successfully deployed from this run."
    fi
else
    echo "Already running the latest build (${CURRENT_COMMIT:0:7})."
fi
