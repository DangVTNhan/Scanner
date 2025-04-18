#!/bin/bash

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file..."
    echo "OPENWEATHER_API_KEY=your_api_key" > .env
    echo ".env file created. Please edit it to add your OpenWeather API key."
    echo ""
fi

# Check if OpenWeather API key is set
if grep -q "your_api_key" .env; then
    echo "WARNING: You need to set your OpenWeather API key in the .env file."
    echo "Please edit the .env file and replace 'your_api_key' with your actual API key."
    echo ""
fi

# Build and start the containers
echo "Building and starting containers..."
docker-compose up -d --build

# Check if containers are running
echo ""
echo "Checking container status..."
docker-compose ps

echo ""
echo "Setup complete!"
echo "The application should be available at: http://localhost:3000"
echo ""
echo "To view logs, run: docker-compose logs -f"
echo "To stop the application, run: docker-compose down"
