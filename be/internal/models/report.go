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
