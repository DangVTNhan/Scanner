package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReportService handles business logic for weather reports
type ReportService struct {
	db             *mongo.Database
	weatherService *openweather.WeatherService
	collection     *mongo.Collection
}

// NewReportService creates a new instance of ReportService
func NewReportService(db *mongo.Database, weatherService *openweather.WeatherService) *ReportService {
	return &ReportService{
		db:             db,
		weatherService: weatherService,
		collection:     db.Collection("reports"),
	}
}

// GenerateReport creates a new weather report
func (s *ReportService) GenerateReport(ctx context.Context, req *models.ReportRequest) (*models.WeatherReport, error) {
	var timestamp time.Time
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	} else {
		timestamp = time.Now()
	}

	var weatherData *openweather.WeatherData
	var err error

	// If timestamp is within the last hour, get current weather
	// Otherwise, get historical weather
	if time.Since(timestamp) < time.Hour {
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

	result, err := s.collection.InsertOne(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	report.ID = result.InsertedID.(primitive.ObjectID)
	return report, nil
}

// GetAllReports retrieves all weather reports
func (s *ReportService) GetAllReports(ctx context.Context) ([]models.WeatherReport, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve reports: %w", err)
	}
	defer cursor.Close(ctx)

	var reports []models.WeatherReport
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, fmt.Errorf("failed to decode reports: %w", err)
	}

	return reports, nil
}

// GetReportByID retrieves a weather report by ID
func (s *ReportService) GetReportByID(ctx context.Context, id string) (*models.WeatherReport, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var report models.WeatherReport
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("report not found")
		}
		return nil, fmt.Errorf("failed to retrieve report: %w", err)
	}

	return &report, nil
}

// CompareReports compares two weather reports
func (s *ReportService) CompareReports(ctx context.Context, req *models.ComparisonRequest) (*models.ComparisonResult, error) {
	report1, err := s.GetReportByID(ctx, req.ReportID1)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve first report: %w", err)
	}

	report2, err := s.GetReportByID(ctx, req.ReportID2)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve second report: %w", err)
	}

	// Calculate deviations
	deviation := models.Deviation{
		Temperature: math.Abs(report2.Temperature - report1.Temperature),
		Pressure:    math.Abs(report2.Pressure - report1.Pressure),
		Humidity:    math.Abs(report2.Humidity - report1.Humidity),
		CloudCover:  math.Abs(report2.CloudCover - report1.CloudCover),
	}

	result := &models.ComparisonResult{
		Report1:   *report1,
		Report2:   *report2,
		Deviation: deviation,
	}

	return result, nil
}
