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
	Reports     []models.WeatherReport `json:"reports"`              // List of reports for the current page
	TotalCount  int                    `json:"totalCount,omitempty"` // Total number of reports (only when not filtered)
	HasMore     bool                   `json:"hasMore"`              // Whether there are more reports to fetch
	CurrentPage int                    `json:"currentPage"`          // Current page number
	FromNumber  int                    `json:"fromNumber"`           // Starting index of the current page
	ToNumber    int                    `json:"toNumber"`             // Ending index of the current page
}
