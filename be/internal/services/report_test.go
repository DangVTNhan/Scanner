package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations of the dependencies

// MockReportRepository is a mock implementation of IReportRepository
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) InsertReport(ctx context.Context, report *models.WeatherReport) (string, error) {
	args := m.Called(ctx, report)
	return args.String(0), args.Error(1)
}

func (m *MockReportRepository) FindAllReports(ctx context.Context) ([]models.WeatherReport, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.WeatherReport), args.Error(1)
}

func (m *MockReportRepository) FindPaginatedReports(ctx context.Context, req *request.PaginatedReportsRequest) (*response.PaginatedReportsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.PaginatedReportsResponse), args.Error(1)
}

func (m *MockReportRepository) FindReportByID(ctx context.Context, id string) (*models.WeatherReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeatherReport), args.Error(1)
}

func (m *MockReportRepository) CountReports(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockWeatherCacheRepository is a mock implementation of IWeatherCacheRepository
type MockWeatherCacheRepository struct {
	mock.Mock
}

func (m *MockWeatherCacheRepository) SaveWeatherCache(ctx context.Context, cache *models.WeatherCache) (string, error) {
	args := m.Called(ctx, cache)
	return args.String(0), args.Error(1)
}

func (m *MockWeatherCacheRepository) FindLatestWeatherCache(ctx context.Context) (*models.WeatherCache, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeatherCache), args.Error(1)
}

func (m *MockWeatherCacheRepository) FindWeatherCacheByTimestamp(ctx context.Context, timestamp time.Time, windowMinutes ...int) (*models.WeatherCache, error) {
	args := m.Called(ctx, timestamp, windowMinutes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeatherCache), args.Error(1)
}

func (m *MockWeatherCacheRepository) DeleteExpiredCaches(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockWeatherService is a mock implementation of IWeatherService
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetCurrentWeather() (*openweather.WeatherData, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*openweather.WeatherData), args.Error(1)
}

func (m *MockWeatherService) GetHistoricalWeather(timestamp time.Time) (*openweather.WeatherData, error) {
	args := m.Called(timestamp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*openweather.WeatherData), args.Error(1)
}

// Test cases

func TestNewReportService(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)

	// Act
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockReportRepo, service.reportRepository)
	assert.Equal(t, mockWeatherCacheRepo, service.weatherCacheRepo)
	assert.Equal(t, mockWeatherService, service.weatherService)
}

func TestGenerateReport_WithTimestamp(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	req := &request.ReportRequest{
		Timestamp: &timestamp,
	}

	// Mock the cache repository to return nil (no cache found)
	mockWeatherCacheRepo.On("FindWeatherCacheByTimestamp", ctx, timestamp, []int{1}).Return(nil, nil)

	// Mock the weather service to return weather data
	weatherData := &openweather.WeatherData{
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
	}
	mockWeatherService.On("GetHistoricalWeather", timestamp).Return(weatherData, nil)

	// Mock the report repository to return an ID
	expectedID := "report123"
	mockReportRepo.On("InsertReport", ctx, mock.AnythingOfType("*models.WeatherReport")).Return(expectedID, nil)

	// Mock the cache repository to save the weather cache
	mockWeatherCacheRepo.On("SaveWeatherCache", ctx, mock.AnythingOfType("*models.WeatherCache")).Return("cache123", nil)

	// Act
	report, err := service.GenerateReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, expectedID, report.ID)
	assert.Equal(t, timestamp, report.Timestamp)
	assert.Equal(t, weatherData.Temperature, report.Temperature)
	assert.Equal(t, weatherData.Pressure, report.Pressure)
	assert.Equal(t, weatherData.Humidity, report.Humidity)
	assert.Equal(t, weatherData.CloudCover, report.CloudCover)

	mockWeatherCacheRepo.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
	mockReportRepo.AssertExpectations(t)
}

func TestGenerateReport_WithoutTimestamp(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	req := &request.ReportRequest{
		Timestamp: nil, // No timestamp provided
	}

	// Mock the cache repository to return nil (no cache found)
	mockWeatherCacheRepo.On("FindWeatherCacheByTimestamp", ctx, mock.AnythingOfType("time.Time"), []int{1}).Return(nil, nil)

	// Mock the weather service to return weather data
	weatherData := &openweather.WeatherData{
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
	}
	mockWeatherService.On("GetCurrentWeather").Return(weatherData, nil)

	// Mock the report repository to return an ID
	expectedID := "report123"
	mockReportRepo.On("InsertReport", ctx, mock.AnythingOfType("*models.WeatherReport")).Return(expectedID, nil)

	// Mock the cache repository to save the weather cache
	mockWeatherCacheRepo.On("SaveWeatherCache", ctx, mock.AnythingOfType("*models.WeatherCache")).Return("cache123", nil)

	// Act
	report, err := service.GenerateReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, expectedID, report.ID)
	assert.Equal(t, weatherData.Temperature, report.Temperature)
	assert.Equal(t, weatherData.Pressure, report.Pressure)
	assert.Equal(t, weatherData.Humidity, report.Humidity)
	assert.Equal(t, weatherData.CloudCover, report.CloudCover)

	mockWeatherCacheRepo.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
	mockReportRepo.AssertExpectations(t)
}

func TestGenerateReport_WithCache(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	req := &request.ReportRequest{
		Timestamp: &timestamp,
	}

	// Mock the cache repository to return a cache
	cacheID := "cache123"
	weatherData := openweather.WeatherData{
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
	}
	cache := &models.WeatherCache{
		ID:          cacheID,
		Timestamp:   timestamp,
		WeatherData: weatherData,
		CreatedAt:   timestamp,
	}
	mockWeatherCacheRepo.On("FindWeatherCacheByTimestamp", ctx, timestamp, []int{1}).Return(cache, nil)

	// Act
	report, err := service.GenerateReport(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, cacheID, report.ID)
	assert.Equal(t, timestamp, report.Timestamp)
	assert.Equal(t, weatherData.Temperature, report.Temperature)
	assert.Equal(t, weatherData.Pressure, report.Pressure)
	assert.Equal(t, weatherData.Humidity, report.Humidity)
	assert.Equal(t, weatherData.CloudCover, report.CloudCover)

	mockWeatherCacheRepo.AssertExpectations(t)
	// Weather service and report repository should not be called
	mockWeatherService.AssertNotCalled(t, "GetCurrentWeather")
	mockWeatherService.AssertNotCalled(t, "GetHistoricalWeather")
	mockReportRepo.AssertNotCalled(t, "InsertReport")
}

func TestGenerateReport_WeatherServiceError(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	req := &request.ReportRequest{
		Timestamp: &timestamp,
	}

	// Mock the cache repository to return nil (no cache found)
	mockWeatherCacheRepo.On("FindWeatherCacheByTimestamp", ctx, timestamp, []int{1}).Return(nil, nil)

	// Mock the weather service to return an error
	expectedErr := errors.New("weather service error")
	mockWeatherService.On("GetHistoricalWeather", timestamp).Return(nil, expectedErr)

	// Act
	report, err := service.GenerateReport(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.Contains(t, err.Error(), "failed to get weather data")
	assert.Contains(t, err.Error(), expectedErr.Error())

	mockWeatherCacheRepo.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
	// Report repository should not be called
	mockReportRepo.AssertNotCalled(t, "InsertReport")
}

func TestGenerateReport_ReportRepositoryError(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	req := &request.ReportRequest{
		Timestamp: &timestamp,
	}

	// Mock the cache repository to return nil (no cache found)
	mockWeatherCacheRepo.On("FindWeatherCacheByTimestamp", ctx, timestamp, []int{1}).Return(nil, nil)

	// Mock the weather service to return weather data
	weatherData := &openweather.WeatherData{
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
	}
	mockWeatherService.On("GetHistoricalWeather", timestamp).Return(weatherData, nil)

	// Mock the report repository to return an error
	expectedErr := errors.New("report repository error")
	mockReportRepo.On("InsertReport", ctx, mock.AnythingOfType("*models.WeatherReport")).Return("", expectedErr)

	// Act
	report, err := service.GenerateReport(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.Contains(t, err.Error(), "failed to save report")
	assert.Contains(t, err.Error(), expectedErr.Error())

	mockWeatherCacheRepo.AssertExpectations(t)
	mockWeatherService.AssertExpectations(t)
	mockReportRepo.AssertExpectations(t)
}

func TestGetAllReports(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	expectedReports := []models.WeatherReport{
		{
			ID:          "report1",
			Timestamp:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
			CreatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:          "report2",
			Timestamp:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			Temperature: 26.5,
			Pressure:    1014.2,
			Humidity:    65.0,
			CloudCover:  35.0,
			CreatedAt:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
		},
	}

	mockReportRepo.On("FindAllReports", ctx).Return(expectedReports, nil)

	// Act
	reports, err := service.GetAllReports(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReports, reports)
	mockReportRepo.AssertExpectations(t)
}

func TestGetPaginatedReports(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	req := &request.PaginatedReportsRequest{
		Limit:  10,
		Offset: 0,
	}

	expectedResponse := &response.PaginatedReportsResponse{
		Reports: []models.WeatherReport{
			{
				ID:          "report2",
				Timestamp:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				Temperature: 26.5,
				Pressure:    1014.2,
				Humidity:    65.0,
				CloudCover:  35.0,
				CreatedAt:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		},
		TotalCount: 10,
	}

	mockReportRepo.On("FindPaginatedReports", ctx, req).Return(expectedResponse, nil)

	// Act
	response, err := service.GetPaginatedReports(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
	mockReportRepo.AssertExpectations(t)
}

func TestGetReportByID(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	reportID := "report1"
	expectedReport := &models.WeatherReport{
		ID:          reportID,
		Timestamp:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	mockReportRepo.On("FindReportByID", ctx, reportID).Return(expectedReport, nil)

	// Act
	report, err := service.GetReportByID(ctx, reportID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReport, report)
	mockReportRepo.AssertExpectations(t)
}

func TestCompareReports(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	req := &request.ComparisonRequest{
		ReportID1: "report1",
		ReportID2: "report2",
	}

	report1 := &models.WeatherReport{
		ID:          "report1",
		Timestamp:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	report2 := &models.WeatherReport{
		ID:          "report2",
		Timestamp:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
		Temperature: 26.5,
		Pressure:    1014.2,
		Humidity:    65.0,
		CloudCover:  35.0,
		CreatedAt:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	mockReportRepo.On("FindReportByID", ctx, req.ReportID1).Return(report1, nil)
	mockReportRepo.On("FindReportByID", ctx, req.ReportID2).Return(report2, nil)

	// Act
	result, err := service.CompareReports(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, *report1, result.Report1)
	assert.Equal(t, *report2, result.Report2)
	assert.Equal(t, 1.0, result.Deviation.Temperature)
	assert.Equal(t, 1.0, result.Deviation.Pressure)
	assert.Equal(t, 5.0, result.Deviation.Humidity)
	assert.Equal(t, 5.0, result.Deviation.CloudCover)
	mockReportRepo.AssertExpectations(t)
}

func TestCompareReports_FirstReportNotFound(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	req := &request.ComparisonRequest{
		ReportID1: "report1",
		ReportID2: "report2",
	}

	expectedErr := errors.New("report not found")
	mockReportRepo.On("FindReportByID", ctx, req.ReportID1).Return(nil, expectedErr)

	// Act
	result, err := service.CompareReports(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve first report")
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockReportRepo.AssertExpectations(t)
}

func TestCompareReports_SecondReportNotFound(t *testing.T) {
	// Arrange
	mockReportRepo := new(MockReportRepository)
	mockWeatherCacheRepo := new(MockWeatherCacheRepository)
	mockWeatherService := new(MockWeatherService)
	service := NewReportService(mockReportRepo, mockWeatherCacheRepo, mockWeatherService)

	ctx := context.Background()
	req := &request.ComparisonRequest{
		ReportID1: "report1",
		ReportID2: "report2",
	}

	report1 := &models.WeatherReport{
		ID:          "report1",
		Timestamp:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	expectedErr := errors.New("report not found")
	mockReportRepo.On("FindReportByID", ctx, req.ReportID1).Return(report1, nil)
	mockReportRepo.On("FindReportByID", ctx, req.ReportID2).Return(nil, expectedErr)

	// Act
	result, err := service.CompareReports(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve second report")
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockReportRepo.AssertExpectations(t)
}
