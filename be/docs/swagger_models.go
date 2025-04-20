package docs

import (
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
)

// This file contains type definitions for Swagger documentation
// These types are not used in the actual code, they're just for documentation

// These are just references to ensure Swagger can find the types
type (
	// WeatherReport is a reference to models.WeatherReport
	WeatherReport struct {
		ID          string    `json:"id" example:"60d21b4667d0d8992e89e9e5"`
		Timestamp   time.Time `json:"timestamp" example:"2023-04-18T12:00:00Z"`
		Temperature float64   `json:"temperature" example:"25.5"` // in Celsius
		Pressure    float64   `json:"pressure" example:"1013.2"`  // in hPa
		Humidity    float64   `json:"humidity" example:"60"`      // in %
		CloudCover  float64   `json:"cloudCover" example:"30"`    // in %
		CreatedAt   time.Time `json:"createdAt" example:"2023-04-18T12:05:00Z"`
	}

	// ReportRequest is a reference to request.ReportRequest
	ReportRequest request.ReportRequest

	// ComparisonRequest is a reference to request.ComparisonRequest
	ComparisonRequest request.ComparisonRequest

	// PaginatedReportsRequest is a reference to request.PaginatedReportsRequest
	PaginatedReportsRequest request.PaginatedReportsRequest

	// BaseResponse is a reference to response.BaseResponse
	BaseResponse response.BaseResponse

	// ComparisonResult is a reference to response.ComparisonResult
	ComparisonResult response.ComparisonResult

	// PaginatedReportsResponse is a reference to response.PaginatedReportsResponse
	PaginatedReportsResponse response.PaginatedReportsResponse
)
