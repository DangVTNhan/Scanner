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

- Go 1.24 or later
- Node.js 18 or later
- MongoDB
- OpenWeather API key
- Docker version 28.0.4 or later

## Setup Instructions

### Using Docker (Recommended)

#### Option 1: Using the setup script

1. Run the setup script:

```bash
./setup.sh
```

1.1 If you met permission denied error, run the following command:

```bash
chmod +x setup.sh
```


2. Edit the `.env` file to add your OpenWeather API key.

3. Access the application at http://localhost:3000

#### Option 2: Manual Docker setup

1. Set your OpenWeather API key in the `.env` file:

```bash
# Edit the .env file
OPENWEATHER_API_KEY=your_api_key
```

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

#### Troubleshooting Docker Issues

If you encounter issues with Docker builds, you can use the provided script:

```bash
./docker-debug.sh
```

This script will help diagnose and fix common Docker build issues.

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

Create a `.env.local` file in the `fe` directory with the following content:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

3. Run the frontend:

```bash
cd fe
npm run dev
```

The frontend will start on http://localhost:3000

## API Endpoints

### Generate Weather Report

```
POST /api/reports
```

### Get All Reports

```
GET /api/reports
```

### Get Report by ID

```
GET /api/reports/{id}
```

### Compare Reports

```
POST /api/reports/compare
```

## Implementation Details

### Backend

- Implemented in Golang with a clean architecture approach
- Uses MongoDB for data storage
- Integrates with OpenWeather API for weather data
- RESTful API design with proper error handling

### Frontend

- Built with Next.js for server-side rendering and routing
- Uses Shadcn UI for a clean, modern interface
- Responsive design for mobile and desktop
- Client-side form validation and error handling
