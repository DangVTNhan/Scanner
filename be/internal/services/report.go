package services

import (
	"context"
	"fmt"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
	"math"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/repository"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
)

// ReportService handles business logic for weather reports
type ReportService struct {
	reportRepository repository.IReportRepository
	weatherCacheRepo repository.IWeatherCacheRepository
	weatherService   openweather.IWeatherService
}

// NewReportService creates a new instance of ReportService
func NewReportService(
	reportRepository repository.IReportRepository,
	weatherCacheRepo repository.IWeatherCacheRepository,
	weatherService openweather.IWeatherService) *ReportService {
	return &ReportService{
		reportRepository: reportRepository,
		weatherCacheRepo: weatherCacheRepo,
		weatherService:   weatherService,
	}
}

// GenerateReport creates a new weather report
func (s *ReportService) GenerateReport(ctx context.Context, req *request.ReportRequest) (*models.WeatherReport, error) {
	var timestamp time.Time
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	} else {
		timestamp = time.Now()
	}

	var weatherData *openweather.WeatherData
	var err error

	// Check if a valid weather cache exists
	cache, err := s.weatherCacheRepo.FindWeatherCacheByTimestamp(ctx, timestamp, 1)
	if err == nil && cache != nil {
		return &models.WeatherReport{
			Timestamp:   timestamp,
			Temperature: cache.WeatherData.Temperature,
			Pressure:    cache.WeatherData.Pressure,
			Humidity:    cache.WeatherData.Humidity,
			CloudCover:  cache.WeatherData.CloudCover,
			CreatedAt:   timestamp,
			ID:          cache.ID,
		}, nil
	}
	// If timestamp is within the last hour, get current weather
	// Otherwise, get historical weather
	if time.Since(timestamp) < 10*time.Minute {
		weatherData, err = s.weatherService.GetCurrentWeather()
	} else {
		weatherData, err = s.weatherService.GetHistoricalWeather(timestamp)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}

	report := &models.WeatherReport{
		Timestamp:   timestamp,
		Temperature: weatherData.Temperature,
		Pressure:    weatherData.Pressure,
		Humidity:    weatherData.Humidity,
		CloudCover:  weatherData.CloudCover,
		CreatedAt:   time.Now(),
	}

	insertedID, err := s.reportRepository.InsertReport(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	// Save the weather data to cache
	cache = &models.WeatherCache{
		Timestamp:   timestamp,
		WeatherData: *weatherData,
		CreatedAt:   time.Now(),
	}

	_, err = s.weatherCacheRepo.SaveWeatherCache(ctx, cache)
	if err != nil {
		// TODO: Handle error (e.g., log it)
	}

	report.ID = insertedID
	return report, nil
}

// GetAllReports retrieves all weather reports (legacy method, kept for backward compatibility)
func (s *ReportService) GetAllReports(ctx context.Context) ([]models.WeatherReport, error) {
	return s.reportRepository.FindAllReports(ctx)
}

// GetPaginatedReports retrieves weather reports with pagination and optional filtering
func (s *ReportService) GetPaginatedReports(ctx context.Context, req *request.PaginatedReportsRequest) (*response.PaginatedReportsResponse, error) {
	return s.reportRepository.FindPaginatedReports(ctx, req)
}

// GetReportByID retrieves a weather report by ID
func (s *ReportService) GetReportByID(ctx context.Context, id string) (*models.WeatherReport, error) {
	return s.reportRepository.FindReportByID(ctx, id)
}

// CompareReports compares two weather reports
func (s *ReportService) CompareReports(ctx context.Context, req *request.ComparisonRequest) (*response.ComparisonResult, error) {
	report1, err := s.GetReportByID(ctx, req.ReportID1)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve first report: %w", err)
	}

	report2, err := s.GetReportByID(ctx, req.ReportID2)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve second report: %w", err)
	}

	// Calculate deviations
	deviation := response.Deviation{
		Temperature: math.Abs(report2.Temperature - report1.Temperature),
		Pressure:    math.Abs(report2.Pressure - report1.Pressure),
		Humidity:    math.Abs(report2.Humidity - report1.Humidity),
		CloudCover:  math.Abs(report2.CloudCover - report1.CloudCover),
	}

	result := &response.ComparisonResult{
		Report1:   *report1,
		Report2:   *report2,
		Deviation: deviation,
	}

	return result, nil
}
