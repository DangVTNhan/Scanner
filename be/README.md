# Changi Airport Weather Report System - Backend

This is the backend for the Changi Airport Weather Report System. It provides APIs for generating weather reports, retrieving historical reports, and comparing reports.

## Project Structure

The project follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) structure:

- `cmd/api`: Main application entry point
- `configs`: Configuration handling
- `internal`: Private application code
  - `handlers`: HTTP handlers
  - `models`: Data models
  - `services`: Business logic
- `pkg`: Public libraries that can be used by external applications
  - `openweather`: OpenWeather API client

## Prerequisites

- Go 1.24 or later
- MongoDB
- OpenWeather API key

## Environment Variables

- `OPENWEATHER_API_KEY`: Your OpenWeather API key (required)
- `MONGO_URI`: MongoDB connection string (default: "mongodb://localhost:27017")
- `DB_NAME`: MongoDB database name (default: "weather_reports")
- `PORT`: Server port (default: "8080")

## Running the Application

1. Set the required environment variables:

```bash
export OPENWEATHER_API_KEY=your_api_key
```

2. Run the application:

```bash
cd cmd/api
go run main.go
```

## API Endpoints

### Generate Weather Report

```
POST /api/reports
```

Request body:
```json
{
  "timestamp": "2023-04-18T12:00:00Z"  // Optional, defaults to current time
}
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

Request body:
```json
{
  "reportId1": "report_id_1",
  "reportId2": "report_id_2"
}
```
