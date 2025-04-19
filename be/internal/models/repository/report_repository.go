package repository

import (
	"context"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
)

// IReportRepository defines the interface for report data access
type IReportRepository interface {
	// InsertReport inserts a new weather report into the database
	InsertReport(ctx context.Context, report *models.WeatherReport) (string, error)

	// FindAllReports retrieves all weather reports
	FindAllReports(ctx context.Context) ([]models.WeatherReport, error)

	// FindPaginatedReports retrieves weather reports with pagination and filtering
	FindPaginatedReports(ctx context.Context, req *request.PaginatedReportsRequest) (*response.PaginatedReportsResponse, error)

	// FindReportByID retrieves a weather report by its ID
	FindReportByID(ctx context.Context, id string) (*models.WeatherReport, error)

	// CountReports counts the total number of reports
	CountReports(ctx context.Context) (int64, error)
}
