package request

import "time"

// ReportRequest represents a request to generate a weather report
type ReportRequest struct {
	Timestamp *time.Time `json:"timestamp"` // Optional: if not provided, current time will be used
}

// ComparisonRequest represents a request to compare two reports
type ComparisonRequest struct {
	ReportID1 string `json:"reportId1"`
	ReportID2 string `json:"reportId2"`
}

// PaginatedReportsRequest represents a request for paginated reports with optional filtering
type PaginatedReportsRequest struct {
	Limit      int       `json:"limit,omitempty"`      // Number of reports to return (default: 10)
	LastID     string    `json:"lastId,omitempty"`     // ID of the last report from the previous page (for cursor-based pagination)
	FromTime   time.Time `json:"fromTime,omitempty"`   // Filter reports from this time
	ToTime     time.Time `json:"toTime,omitempty"`     // Filter reports until this time
	IsFiltered bool      `json:"isFiltered,omitempty"` // Whether filtering is applied
}
