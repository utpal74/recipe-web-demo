#!/bin/bash

# run.sh - Cross-platform shell script to run the Recipes Web API
# Compatible with bash/sh on Linux, macOS, and Windows (Git Bash/WSL)

# Color output for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Trim whitespace function
trim() {
    local var="$@"
    var="${var#"${var%%[![:space:]]*}"}"   # remove leading whitespace
    var="${var%"${var##*[![:space:]]}"}"   # remove trailing whitespace
    echo -n "$var"
}

echo -e "${GREEN}=== Recipes Web API Launcher ===${NC}"

# Check if input is being piped (for automated testing)
if [ ! -t 0 ]; then
    # Read from piped input
    read -r REPO_TYPE_INPUT
    read -r SEED_DATA_INPUT
    REPO_TYPE=$(trim "$REPO_TYPE_INPUT")
    SEED_DATA=$(trim "$SEED_DATA_INPUT")
else
    # Interactive mode
    echo -e "${YELLOW}Available repository types:${NC}"
    echo "  - memory (default: file-based in-memory storage)"
    echo "  - mongo (MongoDB storage with seeding)"
    read -p "Enter REPO_TYPE (default: memory): " REPO_TYPE
    read -p "Enter SEED_DATA (true/false, default: false): " SEED_DATA
fi

REPO_TYPE=${REPO_TYPE:-memory}
SEED_DATA=${SEED_DATA:-false}
export REPO_TYPE
export SEED_DATA

# Validate REPO_TYPE
if [ "$REPO_TYPE" != "memory" ] && [ "$REPO_TYPE" != "mongo" ]; then
    echo -e "${RED}Error: Invalid REPO_TYPE. Must be 'memory' or 'mongo'${NC}"
    echo -e "${RED}Got: '$REPO_TYPE'${NC}"
    exit 1
fi

# Validate SEED_DATA
if [ "$SEED_DATA" != "true" ] && [ "$SEED_DATA" != "false" ]; then
    echo -e "${RED}Error: Invalid SEED_DATA. Must be 'true' or 'false'${NC}"
    echo -e "${RED}Got: '$SEED_DATA'${NC}"
    exit 1
fi

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

# Display configuration
echo ""
echo -e "${GREEN}Configuration:${NC}"
echo "  Repository Type: $REPO_TYPE"
echo "  Seed Data: $SEED_DATA"
echo ""

# Start the application
echo -e "${GREEN}Starting Recipes Web API...${NC}"
echo "The server will be available at: http://localhost:8080"
echo ""
echo "API Endpoints:"
echo "  GET    /recipes              - List all recipes"
echo "  GET    /recipes/{id}         - Get recipe by ID"
echo "  POST   /recipes              - Create new recipe"
echo "  PUT    /recipes/{id}         - Update recipe"
echo "  DELETE /recipes/{id}         - Delete recipe"
echo "  GET    /recipes/search?tag=X - Search recipes by tag"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
echo ""

go run ./cmd/main.go
