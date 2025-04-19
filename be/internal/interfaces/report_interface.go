package interfaces

import (
	"context"
	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
)

type IReportService interface {
	GenerateReport(ctx context.Context, req *request.ReportRequest) (*models.WeatherReport, error)
	GetAllReports(ctx context.Context) ([]models.WeatherReport, error)
	GetPaginatedReports(ctx context.Context, req *request.PaginatedReportsRequest) (*response.PaginatedReportsResponse, error)
	GetReportByID(ctx context.Context, id string) (*models.WeatherReport, error)
	CompareReports(ctx context.Context, req *request.ComparisonRequest) (*response.ComparisonResult, error)
}
