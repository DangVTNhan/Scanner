package mongodb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure that MockSingleResult implements ISingleResult

func TestNewMongoWeatherCacheRepository(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	// Act
	repo := NewMongoWeatherCacheRepository(mockDB)

	// Assert
	assert.NotNil(t, repo)
	mockDB.AssertExpectations(t)
}

func TestSaveWeatherCache(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	cache := &models.WeatherCache{
		Timestamp: time.Now(),
		WeatherData: openweather.WeatherData{
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
		},
		CreatedAt: time.Now(),
	}

	objectID := primitive.NewObjectID()
	mockResult := &mongo.InsertOneResult{
		InsertedID: objectID,
	}

	mockCollection.On("InsertOne", ctx, cache, mock.Anything).Return(mockResult, nil)

	// Act
	id, err := repo.SaveWeatherCache(ctx, cache)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, objectID.Hex(), id)
	mockCollection.AssertExpectations(t)
}

func TestSaveWeatherCache_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	cache := &models.WeatherCache{
		Timestamp: time.Now(),
		WeatherData: openweather.WeatherData{
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
		},
		CreatedAt: time.Now(),
	}

	expectedErr := errors.New("database error")
	mockCollection.On("InsertOne", ctx, cache, mock.Anything).Return(nil, expectedErr)

	// Act
	id, err := repo.SaveWeatherCache(ctx, cache)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestFindLatestWeatherCache(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)
	objectID := primitive.NewObjectID()
	id := objectID.Hex()

	ctx := context.Background()
	now := time.Now()
	// Expected result
	expectedCache := &models.WeatherCache{
		ID:        id,
		Timestamp: now,
		WeatherData: openweather.WeatherData{
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
		},
		CreatedAt: now,
	}

	mockSingleResult := NewMockSingleResult(nil, expectedCache)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindLatestWeatherCache(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.Equal(t, expectedCache.ID, cache.ID)
	assert.Equal(t, expectedCache.Timestamp, cache.Timestamp)
	assert.Equal(t, expectedCache.WeatherData.Temperature, cache.WeatherData.Temperature)
	mockCollection.AssertExpectations(t)
}

func TestFindLatestWeatherCache_NoDocuments(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()

	mockSingleResult := NewMockSingleResult(mongo.ErrNoDocuments, nil)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindLatestWeatherCache(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, cache)
	mockCollection.AssertExpectations(t)
}

func TestFindLatestWeatherCache_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	expectedErr := errors.New("database error")

	mockSingleResult := NewMockSingleResult(expectedErr, nil)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindLatestWeatherCache(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cache)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestFindWeatherCacheByTimestamp(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything, mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	timestamp := time.Now()

	// Expected cache
	expectedCache := &models.WeatherCache{
		ID:        "cache123",
		Timestamp: timestamp,
		WeatherData: openweather.WeatherData{
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
		},
		CreatedAt: timestamp,
	}

	mockSingleResult := NewMockSingleResult(nil, expectedCache)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindWeatherCacheByTimestamp(ctx, timestamp)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.Equal(t, expectedCache.ID, cache.ID)
	assert.Equal(t, expectedCache.Timestamp, cache.Timestamp)
	assert.Equal(t, expectedCache.WeatherData.Temperature, cache.WeatherData.Temperature)
	mockCollection.AssertExpectations(t)
}

func TestFindWeatherCacheByTimestamp_CustomWindow(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	timestamp := time.Now()
	window := 5 // Custom window

	// Expected cache
	expectedCache := &models.WeatherCache{
		ID:        "cache123",
		Timestamp: timestamp,
		WeatherData: openweather.WeatherData{
			Temperature: 25.5,
			Pressure:    1013.2,
			Humidity:    60.0,
			CloudCover:  30.0,
		},
		CreatedAt: timestamp,
	}

	mockSingleResult := NewMockSingleResult(nil, expectedCache)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindWeatherCacheByTimestamp(ctx, timestamp, window)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.Equal(t, expectedCache.ID, cache.ID)
	assert.Equal(t, expectedCache.Timestamp, cache.Timestamp)
	assert.Equal(t, expectedCache.WeatherData.Temperature, cache.WeatherData.Temperature)
	mockCollection.AssertExpectations(t)
}

func TestFindWeatherCacheByTimestamp_NoDocuments(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	timestamp := time.Now()

	mockSingleResult := NewMockSingleResult(mongo.ErrNoDocuments, nil)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindWeatherCacheByTimestamp(ctx, timestamp)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, cache)
	mockCollection.AssertExpectations(t)
}

func TestFindWeatherCacheByTimestamp_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	timestamp := time.Now()
	expectedErr := errors.New("database error")

	mockSingleResult := NewMockSingleResult(expectedErr, nil)
	mockCollection.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSingleResult)

	// Act
	cache, err := repo.FindWeatherCacheByTimestamp(ctx, timestamp)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cache)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}

func TestDeleteExpiredCaches(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	mockResult := &mongo.DeleteResult{
		DeletedCount: 5,
	}

	mockCollection.On("DeleteMany", ctx, mock.Anything, mock.Anything).Return(mockResult, nil)

	// Act
	err := repo.DeleteExpiredCaches(ctx)

	// Assert
	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestDeleteExpiredCaches_Error(t *testing.T) {
	// Arrange
	mockDB := new(MockDatabase)
	mockCollection := new(MockCollection)
	mockDB.On("Collection", "weather_cache", mock.Anything).Return(mockCollection)

	repo := NewMongoWeatherCacheRepository(mockDB)

	ctx := context.Background()
	expectedErr := errors.New("database error")

	mockCollection.On("DeleteMany", ctx, mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Act
	err := repo.DeleteExpiredCaches(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockCollection.AssertExpectations(t)
}
