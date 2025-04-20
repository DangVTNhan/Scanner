package response

import "github.com/DangVTNhan/Scanner/be/internal/models"

// ComparisonResult represents the result of comparing two reports
type ComparisonResult struct {
	Report1   models.WeatherReport `json:"report1"`
	Report2   models.WeatherReport `json:"report2"`
	Deviation Deviation            `json:"deviation"`
}

// Deviation represents the differences between two reports
type Deviation struct {
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
	CloudCover  float64 `json:"cloudCover"`
}

// PaginatedReportsResponse represents a paginated response of weather reports
type PaginatedReportsResponse struct {
	Reports    []models.WeatherReport `json:"reports"`    // List of reports for the current page
	TotalCount int                    `json:"totalCount"` // Total number of reports (for calculating total pages)
}
