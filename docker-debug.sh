#!/bin/bash

# This script helps debug Docker build issues

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Error: Docker is not running or not installed."
  exit 1
fi

# Clean up any previous builds
echo "Cleaning up previous builds..."
docker-compose down -v
docker system prune -f

# Build the frontend image separately with verbose output
echo "Building frontend image..."
cd fe
docker build -t scanner-frontend . 2>&1 | tee ../frontend-build.log
cd ..

# Check if the build was successful
if [ $? -eq 0 ]; then
  echo "Frontend build successful!"
  echo "You can now run 'docker-compose up -d' to start all services."
else
  echo "Frontend build failed. Check frontend-build.log for details."
  echo "Common issues:"
  echo "1. Missing next.config.js file"
  echo "2. Incompatible dependencies"
  echo "3. Build errors in the Next.js application"
fi
