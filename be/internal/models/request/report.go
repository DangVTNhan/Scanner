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

// SortOrder represents the sort order (ascending or descending)
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// PaginatedReportsRequest represents a request for paginated reports with optional filtering
type PaginatedReportsRequest struct {
	Limit      int       `json:"limit,omitempty"`      // Number of reports to return (default: 10)
	Offset     int       `json:"offset,omitempty"`     // Number of reports to skip (for pagination)
	FromTime   time.Time `json:"fromTime,omitempty"`   // Filter reports from this time
	ToTime     time.Time `json:"toTime,omitempty"`     // Filter reports until this time
	IsFiltered bool      `json:"isFiltered,omitempty"` // Whether filtering is applied
	SortBy     string    `json:"sortBy,omitempty"`     // Field to sort by (default: "timestamp")
	SortOrder  SortOrder `json:"sortOrder,omitempty"`  // Sort order (default: "desc")
}
