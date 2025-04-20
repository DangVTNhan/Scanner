package mongodb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewMongoReportRepository(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	// Act
	repo := NewMongoReportRepository(mockDB)

	// Assert
	assert.NotNil(t, repo)
	mockDB.AssertExpectations(t)
}

func TestInsertReport(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	report := &models.WeatherReport{
		Timestamp:   time.Now(),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Now(),
	}

	objectID := primitive.NewObjectID()
	mockResult := &mongo.InsertOneResult{
		InsertedID: objectID,
	}

	mockCollection.On("InsertOne", ctx, report, mock.Anything).Return(mockResult, nil)

	// Act
	id, err := repo.InsertReport(ctx, report)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, objectID.Hex(), id)
	mockCollection.AssertExpectations(t)
}

func TestInsertReport_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	report := &models.WeatherReport{
		Timestamp:   time.Now(),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Now(),
	}

	expectedErr := errors.New("database error")
	mockCollection.On("InsertOne", ctx, report, mock.Anything).Return(nil, expectedErr)

	// Act
	id, err := repo.InsertReport(ctx, report)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestFindAllReports(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	now := time.Now()

	// Expected reports
	expectedReports := []models.WeatherReport{
		{
			ID:          "report1",
			Timestamp:   now,
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
			CreatedAt:   now,
		},
		{
			ID:          "report2",
			Timestamp:   now.Add(1 * time.Hour),
			Temperature: 26.5,
			Pressure:    1014.2,
			Humidity:    65.0,
			CloudCover:  35.0,
			CreatedAt:   now.Add(1 * time.Hour),
		},
	}

	mockCursor := NewMockCursor(expectedReports)
	mockCursor.On("All", ctx, mock.AnythingOfType("*[]models.WeatherReport")).Return(nil)
	mockCursor.On("Close", ctx).Return(nil)

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)

	// Act
	reports, err := repo.FindAllReports(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(expectedReports), len(reports))
	assert.Equal(t, expectedReports[0].ID, reports[0].ID)
	assert.Equal(t, expectedReports[1].ID, reports[1].ID)
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestFindAllReports_FindError(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	expectedErr := errors.New("database error")

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Act
	reports, err := repo.FindAllReports(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, reports)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestFindAllReports_DecodeError(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	expectedErr := errors.New("decode error")

	mockCursor := new(MockCursor)
	mockCursor.On("All", ctx, mock.AnythingOfType("*[]models.WeatherReport")).Return(expectedErr)
	mockCursor.On("Close", ctx).Return(nil)

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)

	// Act
	reports, err := repo.FindAllReports(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, reports)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestFindPaginatedReports(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	now := time.Now()
	req := &request.PaginatedReportsRequest{
		Limit:  10,
		LastID: "report1",
	}

	// Expected reports
	expectedReports := []models.WeatherReport{
		{
			ID:          "report2",
			Timestamp:   now,
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
			CreatedAt:   now,
		},
		{
			ID:          "report3",
			Timestamp:   now.Add(1 * time.Hour),
			Temperature: 26.5,
			Pressure:    1014.2,
			Humidity:    65.0,
			CloudCover:  35.0,
			CreatedAt:   now.Add(1 * time.Hour),
		},
	}

	mockCursor := NewMockCursor(expectedReports)
	mockCursor.On("All", ctx, mock.AnythingOfType("*[]models.WeatherReport")).Return(nil)
	mockCursor.On("Close", ctx).Return(nil)

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)
	mockCollection.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(10), nil)

	// Act
	response, err := repo.FindPaginatedReports(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, len(expectedReports), len(response.Reports))
	assert.Equal(t, expectedReports[0].ID, response.Reports[0].ID)
	assert.Equal(t, expectedReports[1].ID, response.Reports[1].ID)
	assert.Equal(t, 10, response.TotalCount)
	assert.Equal(t, 2, response.CurrentPage)
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestFindPaginatedReports_WithFilters(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	now := time.Now()
	fromTime := now.Add(-24 * time.Hour)
	toTime := now
	req := &request.PaginatedReportsRequest{
		Limit:      10,
		FromTime:   fromTime,
		ToTime:     toTime,
		IsFiltered: true,
	}

	// Expected reports
	expectedReports := []models.WeatherReport{
		{
			ID:          "report1",
			Timestamp:   now.Add(-12 * time.Hour),
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
			CreatedAt:   now.Add(-12 * time.Hour),
		},
	}

	mockCursor := NewMockCursor(expectedReports)
	mockCursor.On("All", ctx, mock.AnythingOfType("*[]models.WeatherReport")).Return(nil)
	mockCursor.On("Close", ctx).Return(nil)

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)

	// Act
	response, err := repo.FindPaginatedReports(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, len(expectedReports), len(response.Reports))
	assert.Equal(t, expectedReports[0].ID, response.Reports[0].ID)
	assert.Equal(t, 0, response.TotalCount) // TotalCount should be 0 for filtered results
	assert.Equal(t, 1, response.CurrentPage)
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestFindPaginatedReports_FindError(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	req := &request.PaginatedReportsRequest{
		Limit: 10,
	}
	expectedErr := errors.New("database error")

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Act
	response, err := repo.FindPaginatedReports(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestFindPaginatedReports_DecodeError(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	req := &request.PaginatedReportsRequest{
		Limit: 10,
	}
	expectedErr := errors.New("decode error")

	mockCursor := new(MockCursor)
	mockCursor.On("All", ctx, mock.AnythingOfType("*[]models.WeatherReport")).Return(expectedErr)
	mockCursor.On("Close", ctx).Return(nil)

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)

	// Act
	response, err := repo.FindPaginatedReports(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestFindPaginatedReports_CountError(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	req := &request.PaginatedReportsRequest{
		Limit: 10,
	}
	expectedErr := errors.New("count error")

	// Expected reports
	expectedReports := []models.WeatherReport{
		{
			ID:          "report1",
			Timestamp:   time.Now(),
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
			CreatedAt:   time.Now(),
		},
	}

	mockCursor := NewMockCursor(expectedReports)
	mockCursor.On("All", ctx, mock.AnythingOfType("*[]models.WeatherReport")).Return(nil)
	mockCursor.On("Close", ctx).Return(nil)

	mockCollection.On("Find", ctx, mock.Anything, mock.Anything).Return(mockCursor, nil)
	mockCollection.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(0), expectedErr)

	// Act
	response, err := repo.FindPaginatedReports(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestFindReportByID_ValidObjectID(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	objectID := primitive.NewObjectID()
	id := objectID.Hex()

	expectedReport := &models.WeatherReport{
		ID:          id,
		Timestamp:   time.Now(),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Now(),
	}

	mockSingleResult := NewMockSingleResult(nil, expectedReport)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	report, err := repo.FindReportByID(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, expectedReport.ID, report.ID)
	assert.Equal(t, expectedReport.Temperature, report.Temperature)
	mockCollection.AssertExpectations(t)
}

func TestFindReportByID_InvalidObjectID(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	id := "invalid-id"

	expectedReport := &models.WeatherReport{
		ID:          id,
		Timestamp:   time.Now(),
		Temperature: 25.5,
		Pressure:    1013.2,
		Humidity:    60.0,
		CloudCover:  30.0,
		CreatedAt:   time.Now(),
	}

	mockSingleResult := NewMockSingleResult(nil, expectedReport)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	report, err := repo.FindReportByID(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, expectedReport.ID, report.ID)
	assert.Equal(t, expectedReport.Temperature, report.Temperature)
	mockCollection.AssertExpectations(t)
}

func TestFindReportByID_NotFound(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	objectID := primitive.NewObjectID()
	id := objectID.Hex()

	mockSingleResult := NewMockSingleResult(mongo.ErrNoDocuments, nil)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	report, err := repo.FindReportByID(ctx, id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.Contains(t, err.Error(), "report not found")
	mockCollection.AssertExpectations(t)
}

func TestFindReportByID_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	objectID := primitive.NewObjectID()
	id := objectID.Hex()
	expectedErr := errors.New("database error")

	mockSingleResult := NewMockSingleResult(expectedErr, nil)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	report, err := repo.FindReportByID(ctx, id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestCountReports(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	expectedCount := int64(10)

	mockCollection.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(expectedCount, nil)

	// Act
	count, err := repo.CountReports(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockCollection.AssertExpectations(t)
}

func TestCountReports_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "reports", mock.Anything).Return(mockCollection)

	repo := NewMongoReportRepository(mockDB)

	ctx := context.Background()
	expectedErr := errors.New("database error")

	mockCollection.On("CountDocuments", ctx, mock.Anything, mock.Anything).Return(int64(0), expectedErr)

	// Act
	count, err := repo.CountReports(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}
