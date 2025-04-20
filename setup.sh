#!/bin/bash

# Color codes for output formatting
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
NC="\033[0m" # No Color

# Required versions
REQUIRED_GO_VERSION="1.24"
REQUIRED_NODE_VERSION="18"
REQUIRED_DOCKER_VERSION="28.0.4"

# Variables to track dependency status
GO_OK=false
NODE_OK=false
DOCKER_OK=false

# Function to check if Go is installed with the correct version
check_go() {
    echo "Checking Go installation..."
    if ! command -v go &>/dev/null; then
        echo -e "${RED}Go is not installed.${NC}"
        echo "Please install Go version $REQUIRED_GO_VERSION or later from: https://go.dev/doc/install"
        return 1
    fi

    # Get Go version
    GO_VERSION=$(go version | grep -oE 'go[0-9]+(\.[0-9]+)+' | cut -c 3-)
    GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
    GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
    REQ_MAJOR=$(echo $REQUIRED_GO_VERSION | cut -d. -f1)
    REQ_MINOR=$(echo $REQUIRED_GO_VERSION | cut -d. -f2)

    # Compare versions
    if [ "$GO_MAJOR" -gt "$REQ_MAJOR" ] || ([ "$GO_MAJOR" -eq "$REQ_MAJOR" ] && [ "$GO_MINOR" -ge "$REQ_MINOR" ]); then
        echo -e "${GREEN}Go version $GO_VERSION is installed (required: $REQUIRED_GO_VERSION).${NC}"
        GO_OK=true
        return 0
    else
        echo -e "${RED}Go version $GO_VERSION is installed, but version $REQUIRED_GO_VERSION or later is required.${NC}"
        echo "Please update Go from: https://go.dev/doc/install"
        return 1
    fi
}

# Function to check if Node.js is installed with the correct version
check_node() {
    echo "Checking Node.js installation..."
    if ! command -v node &>/dev/null; then
        echo -e "${RED}Node.js is not installed.${NC}"
        echo "Please install Node.js version $REQUIRED_NODE_VERSION or later from: https://nodejs.org/en/download/"
        return 1
    fi

    # Get Node.js version
    NODE_VERSION=$(node -v | cut -c 2-)
    NODE_MAJOR=$(echo $NODE_VERSION | cut -d. -f1)
    REQ_NODE_MAJOR=$(echo $REQUIRED_NODE_VERSION | cut -d. -f1)

    # Compare versions
    if [ "$NODE_MAJOR" -ge "$REQ_NODE_MAJOR" ]; then
        echo -e "${GREEN}Node.js version $NODE_VERSION is installed (required: $REQUIRED_NODE_VERSION or later).${NC}"
        NODE_OK=true
        return 0
    else
        echo -e "${RED}Node.js version $NODE_VERSION is installed, but version $REQUIRED_NODE_VERSION or later is required.${NC}"
        echo "Please update Node.js from: https://nodejs.org/en/download/"
        return 1
    fi
}

# Function to check if Docker is installed with the correct version
check_docker() {
    echo "Checking Docker installation..."
    if ! command -v docker &>/dev/null; then
        echo -e "${RED}Docker is not installed.${NC}"
        echo "Please install Docker version $REQUIRED_DOCKER_VERSION or later from: https://docs.docker.com/get-docker/"
        return 1
    fi

    # Check if Docker daemon is running
    if ! docker info &>/dev/null; then
        echo -e "${RED}Docker daemon is not running.${NC}"
        echo "Please start the Docker daemon and try again."
        return 1
    fi

    # Get Docker version
    DOCKER_VERSION=$(docker version --format '{{.Server.Version}}' 2>/dev/null)
    if [ -z "$DOCKER_VERSION" ]; then
        DOCKER_VERSION=$(docker version | grep -i "version" | head -n 1 | awk '{print $2}')
    fi

    # Extract major, minor, patch versions
    DOCKER_MAJOR=$(echo $DOCKER_VERSION | cut -d. -f1)
    DOCKER_MINOR=$(echo $DOCKER_VERSION | cut -d. -f2)
    DOCKER_PATCH=$(echo $DOCKER_VERSION | cut -d. -f3)
    REQ_DOCKER_MAJOR=$(echo $REQUIRED_DOCKER_VERSION | cut -d. -f1)
    REQ_DOCKER_MINOR=$(echo $REQUIRED_DOCKER_VERSION | cut -d. -f2)
    REQ_DOCKER_PATCH=$(echo $REQUIRED_DOCKER_VERSION | cut -d. -f3)

    # Compare versions
    if [ "$DOCKER_MAJOR" -gt "$REQ_DOCKER_MAJOR" ] ||
        ([ "$DOCKER_MAJOR" -eq "$REQ_DOCKER_MAJOR" ] && [ "$DOCKER_MINOR" -gt "$REQ_DOCKER_MINOR" ]) ||
        ([ "$DOCKER_MAJOR" -eq "$REQ_DOCKER_MAJOR" ] && [ "$DOCKER_MINOR" -eq "$REQ_DOCKER_MINOR" ] && [ "$DOCKER_PATCH" -ge "$REQ_DOCKER_PATCH" ]); then
        echo -e "${GREEN}Docker version $DOCKER_VERSION is installed (required: $REQUIRED_DOCKER_VERSION).${NC}"
        DOCKER_OK=true
        return 0
    else
        echo -e "${RED}Docker version $DOCKER_VERSION is installed, but version $REQUIRED_DOCKER_VERSION or later is required.${NC}"
        echo "Please update Docker from: https://docs.docker.com/get-docker/"
        return 1
    fi
}

# Function to check if docker-compose is available
check_docker_compose() {
    echo "Checking Docker Compose installation..."
    if command -v docker-compose &>/dev/null; then
        echo -e "${GREEN}Docker Compose is installed.${NC}"
        return 0
    elif docker compose version &>/dev/null; then
        echo -e "${GREEN}Docker Compose plugin is installed.${NC}"
        # Create an alias for docker-compose if it doesn't exist
        if ! command -v docker-compose &>/dev/null; then
            echo "Creating alias for 'docker compose' as 'docker-compose'..."
            alias docker-compose="docker compose"
        fi
        return 0
    else
        echo -e "${RED}Docker Compose is not installed.${NC}"
        echo "Please install Docker Compose from: https://docs.docker.com/compose/install/"
        return 1
    fi
}

# Function to display a summary of dependency checks
display_summary() {
    echo ""
    echo "=== Dependency Check Summary ==="

    if $GO_OK; then
        echo -e "${GREEN}✓ Go${NC}"
    else
        echo -e "${RED}✗ Go${NC}"
    fi

    if $NODE_OK; then
        echo -e "${GREEN}✓ Node.js${NC}"
    else
        echo -e "${RED}✗ Node.js${NC}"
    fi

    if $DOCKER_OK; then
        echo -e "${GREEN}✓ Docker${NC}"
    else
        echo -e "${RED}✗ Docker${NC}"
    fi

    echo "==========================="
}

# Check dependencies
echo "Checking dependencies..."
check_go
check_node
check_docker
check_docker_compose
display_summary

# Prompt user to continue if any dependency check failed
if ! $GO_OK || ! $NODE_OK || ! $DOCKER_OK; then
    echo -e "${YELLOW}Some dependency checks failed.${NC}"
    read -p "Do you want to continue anyway? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup aborted."
        exit 1
    fi
    echo "Continuing with setup..."
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file..."
    echo "OPENWEATHER_API_KEY=your_api_key" >.env
    echo ".env file created. Please edit it to add your OpenWeather API key."
    echo ""
fi

# Check if OpenWeather API key is set
if grep -q "your_api_key" .env; then
    echo -e "${YELLOW}WARNING: You need to set your OpenWeather API key in the .env file.${NC}"
    echo "Please edit the .env file and replace 'your_api_key' with your actual API key."
    echo ""
fi

# Build and start the containers
echo "Building and starting containers..."

# Use docker compose or docker-compose based on availability
if command -v docker-compose &>/dev/null; then
    docker-compose up -d --build

    # Check if containers are running
    echo ""
    echo "Checking container status..."
    docker-compose ps
else
    docker compose up -d --build

    # Check if containers are running
    echo ""
    echo "Checking container status..."
    docker compose ps
fi

echo ""
echo -e "${GREEN}Setup complete!${NC}"
echo "The application should be available at: http://localhost:3000"
echo ""

# Show appropriate commands based on which docker compose command is available
if command -v docker-compose &>/dev/null; then
    echo "To view logs, run: docker-compose logs -f"
    echo "To stop the application, run: docker-compose down"
else
    echo "To view logs, run: docker compose logs -f"
    echo "To stop the application, run: docker compose down"
fi
