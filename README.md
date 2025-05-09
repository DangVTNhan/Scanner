# Changi Airport Weather Report System (Scanner)

A full-stack application for generating, storing, and comparing weather reports for Changi Airport using the OpenWeatherMap API.

## Project Structure

The project is organized into two main directories:

- `be/`: Backend implemented in Golang following the golang-standards/project-layout structure
- `fe/`: Frontend implemented with Next.js and Shadcn UI

## Features

- Generate weather reports for Changi Airport (current or historical)
- View a history of all generated reports
- Compare two reports to see deviations in weather parameters
- Store reports in MongoDB for future reference

## Prerequisites

- Go 1.24 or later (install at: https://go.dev/doc/install)
- Node.js 18 or later (install at:https://nodejs.org/en/download/)
- MongoDB (recommend using Docker)
- OpenWeather API key (Register One Call API 3.0: https://openweathermap.org/api/one-call-3)
- Docker version 28.0.4 or later

## Setup Instructions

### Using Docker (Recommended)

#### Option 1: Using the setup script

1. Create a `.env` file in the root directory with the following content:

```
OPENWEATHER_API_KEY=your_api_key
ENVIRONMENT=dev
```

Change `your_api_key` to your actual OpenWeather API key.

If you need openweather api key, you can get it from [here](https://openweathermap.org/api/one-call-3) it required credit card information. Or you can DM me at (nhan.dangviettrung@gmail.com) to get a free api key.


2. Run the setup script:

```bash
./setup.sh
```

If you met permission denied error, run the following command:

```bash
chmod +x setup.sh
```

3. Access the application at http://localhost:3000

#### Option 2: Manual Docker setup

1. Set your OpenWeather API key in the `.env` file:

```bash
# Edit the .env file
OPENWEATHER_API_KEY=your_api_key
```

By default, docker compose will load the `.env` file automatically for environment variables.

2. Build and run the application using Docker Compose:

```bash
docker-compose up -d --build
```

This will start:

- MongoDB on port 27017
- Backend on port 8080
- Frontend on port 3000

Access the application at http://localhost:3000

#### Docker Commands

```bash
# View logs
docker-compose logs -f

# Stop the application
docker-compose down

# Stop the application and remove volumes
docker-compose down -v
```

### Manual Setup

#### Backend

1. Set up environment variables:

```bash
export OPENWEATHER_API_KEY=your_api_key
```

2. Start MongoDB (if not already running):

```bash
mongod --dbpath /path/to/data/db
```

3. Run the backend:

```bash
cd be/cmd/api
go run main.go
```

The backend will start on http://localhost:8080

#### Frontend

1. Install dependencies:

```bash
cd fe
npm install
```

2. Set up environment variables:

Create a `.env.local` file in the `fe` directory with the following content: (if not provided it will use the default value: http://localhost:8080/api)

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

3. Run the frontend:

```bash
cd fe
npm run dev
```

The frontend will start on http://localhost:3000

## Testing

### Backend Tests

Run the backend tests:

```bash
make test-be
# or
cd be && go test ./...
```

### Frontend Tests

Run the frontend tests:

```bash
make test-fe
# or
cd fe && npm test
```

Generate frontend test coverage report:

```bash
make test-fe-coverage
# or
cd fe && npm run test:coverage 
```

## API Documentation

### Swagger UI

The API documentation is available through Swagger UI in development and staging environments:

```
http://localhost:8080/swagger/index.html
```

To enable Swagger UI, make sure the `ENVIRONMENT` variable is set to either `dev` or `stg` in your `.env` file:

```
ENVIRONMENT=dev
```

To generate or update the Swagger documentation, run:

```bash
make generate-swagger
```

If you encounter any issues with Swagger generation, try:

1. Make sure all model types referenced in your Swagger annotations are properly imported
2. Check that the Go module path in your go.mod file matches the import paths in your code
3. If you get redeclaration errors, check for duplicate type definitions

### API Endpoints

#### Generate Weather Report

```
POST /api/reports
```

#### Get All Reports

```
GET /api/reports
```

#### Get Report by ID

```
GET /api/reports/{id}
```

#### Compare Reports

```
POST /api/reports/compare
```

## Implementation Details

- Docker support for easy deployment and scaling
- Makefile and setup script for simplifying common development tasks

### Backend

- Implemented in Golang with a clean architecture approach
- Uses MongoDB for data storage
- Integrates with OpenWeather API for weather data
- Caching mechanism for weather data to reduce API calls (singleflight)
- Dependency injection for better modularity and testability
- RESTful API design with proper error handling
- CORS middleware for handling cross-origin requests
- Swagger documentation for API endpoints
- Extensive test suite using Go's testing package
- Environment variables for configuration (e.g., API keys, database connection strings)

### Frontend

- Built with Next.js for server-side rendering and routing
- Uses Shadcn UI for a clean, modern interface
- Responsive design for mobile and desktop
- Client-side form validation and error handling
- Comprehensive test suite using Jest and React Testing Library
- Environment variables for API configuration
