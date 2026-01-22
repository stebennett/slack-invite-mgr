#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a variable is set
check_env_var() {
    if [ -z "${!1}" ]; then
        echo -e "${RED}Error: $1 is not set${NC}"
        echo "Please set it in your environment or .env file"
        exit 1
    fi
}

# Function to check if a file exists
check_file() {
    if [ ! -f "$1" ]; then
        echo -e "${RED}Error: File $1 does not exist${NC}"
        exit 1
    fi
}

# Check for required environment variables
echo -e "${YELLOW}Checking environment variables...${NC}"
check_env_var "GOOGLE_CREDENTIALS_FILE"
check_env_var "GOOGLE_SPREADSHEET_ID"
check_env_var "GOOGLE_SHEET_NAME"
check_env_var "APPRISE_URL"
check_env_var "GITHUB_USERNAME"

# Check if credentials file exists
echo -e "${YELLOW}Checking credentials file...${NC}"
check_file "$GOOGLE_CREDENTIALS_FILE"

# Create data directory if it doesn't exist
echo -e "${YELLOW}Setting up data directory...${NC}"
mkdir -p data

# Copy credentials file to data directory
echo -e "${YELLOW}Copying credentials file...${NC}"
cp "$GOOGLE_CREDENTIALS_FILE" data/credentials.json

# Stop any existing sheets container
echo -e "${YELLOW}Stopping any existing sheets container...${NC}"
docker stop slack-invite-mgr-sheets 2>/dev/null || true
docker rm slack-invite-mgr-sheets 2>/dev/null || true

# Start the sheets service
echo -e "${YELLOW}Starting sheets service...${NC}"
docker run \
    --name slack-invite-mgr-sheets \
    --rm \
    -v "$(pwd)/data/credentials.json:/app/credentials/credentials.json" \
    -e GOOGLE_CREDENTIALS_FILE=/app/credentials/credentials.json \
    -e GOOGLE_SPREADSHEET_ID="$GOOGLE_SPREADSHEET_ID" \
    -e GOOGLE_SHEET_NAME="$GOOGLE_SHEET_NAME" \
    -e APPRISE_URL="$APPRISE_URL" \
    ghcr.io/$GITHUB_USERNAME/slack-invite-mgr-sheets:latest

# Check if the container started successfully
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Sheets service started successfully!${NC}"
    echo -e "${GREEN}Container logs:${NC}"
    docker logs -f slack-invite-mgr-sheets
else
    echo -e "${RED}Failed to start sheets service${NC}"
    exit 1
fi 