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
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins for CORS (default: "http://localhost:3000,http://frontend:3000,http://host.docker.internal:3000,*")

## CORS Configuration

The application includes CORS middleware to handle Cross-Origin Resource Sharing. By default, the following origins are allowed:

- http://localhost:3000
- http://frontend:3000
- http://host.docker.internal:3000
- * (wildcard for development)

You can customize the allowed origins by setting the `CORS_ALLOWED_ORIGINS` environment variable with a comma-separated list of origins:

```bash
export CORS_ALLOWED_ORIGINS="https://your-domain.com,https://another-domain.com"
```

## Running the Application

1. Set the required environment variables:

```bash
export OPENWEATHER_API_KEY=your_api_key
```

1.1 If you need openweather api key, you can get it from [here](https://openweathermap.org/api/one-call-3) it required credit card information. Or you can DM me at (nhan.dangviettrung@gmail.com) to get a free api key.

2. Run the application:

```bash
cd cmd/api
go run main.go
```

## Testing

To run the tests:

```bash
go test ./...
```

### Testing CORS Middleware

The CORS middleware tests use environment variables to control the allowed origins. When running the tests, the `CORS_ALLOWED_ORIGINS` environment variable is temporarily set to specific values for each test case, and then restored to its original value after the test completes.

This approach allows testing different CORS configurations without modifying the actual code or creating complex mocks.

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
