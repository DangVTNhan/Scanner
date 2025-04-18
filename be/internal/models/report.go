package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WeatherReport represents a weather report for Changi Airport
type WeatherReport struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Temperature float64            `json:"temperature" bson:"temperature"` // in Celsius
	Pressure    float64            `json:"pressure" bson:"pressure"`       // in hPa
	Humidity    float64            `json:"humidity" bson:"humidity"`       // in %
	CloudCover  float64            `json:"cloudCover" bson:"cloudCover"`   // in %
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

// ReportRequest represents a request to generate a weather report
type ReportRequest struct {
	Timestamp *time.Time `json:"timestamp"` // Optional: if not provided, current time will be used
}

// ComparisonRequest represents a request to compare two reports
type ComparisonRequest struct {
	ReportID1 string `json:"reportId1"`
	ReportID2 string `json:"reportId2"`
}

// ComparisonResult represents the result of comparing two reports
type ComparisonResult struct {
	Report1   WeatherReport `json:"report1"`
	Report2   WeatherReport `json:"report2"`
	Deviation Deviation     `json:"deviation"`
}

// Deviation represents the differences between two reports
type Deviation struct {
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
	CloudCover  float64 `json:"cloudCover"`
}

// PaginatedReportsRequest represents a request for paginated reports with optional filtering
type PaginatedReportsRequest struct {
	Limit     int       `json:"limit,omitempty"`     // Number of reports to return (default: 10)
	LastID    string    `json:"lastId,omitempty"`    // ID of the last report from the previous page (for cursor-based pagination)
	FromTime  time.Time `json:"fromTime,omitempty"`  // Filter reports from this time
	ToTime    time.Time `json:"toTime,omitempty"`    // Filter reports until this time
	IsFiltered bool      `json:"isFiltered,omitempty"` // Whether filtering is applied
}

// PaginatedReportsResponse represents a paginated response of weather reports
type PaginatedReportsResponse struct {
	Reports     []WeatherReport `json:"reports"`               // List of reports for the current page
	TotalCount  int             `json:"totalCount,omitempty"`  // Total number of reports (only when not filtered)
	HasMore     bool            `json:"hasMore"`              // Whether there are more reports to fetch
	CurrentPage int             `json:"currentPage"`          // Current page number
	FromNumber  int             `json:"fromNumber"`           // Starting index of the current page
	ToNumber    int             `json:"toNumber"`             // Ending index of the current page
}
