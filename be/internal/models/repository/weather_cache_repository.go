package repository

import (
	"context"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
)

// IWeatherCacheRepository defines the interface for weather cache data access
type IWeatherCacheRepository interface {
	// SaveWeatherCache saves a weather data cache entry
	SaveWeatherCache(ctx context.Context, cache *models.WeatherCache) (string, error)

	// FindLatestWeatherCache retrieves the latest valid weather cache entry
	FindLatestWeatherCache(ctx context.Context) (*models.WeatherCache, error)

	// FindWeatherCacheByTimestamp retrieves a weather cache entry by timestamp within a time window
	FindWeatherCacheByTimestamp(ctx context.Context, timestamp time.Time, windowMinutes ...int) (*models.WeatherCache, error)

	// DeleteExpiredCaches removes expired cache entries
	DeleteExpiredCaches(ctx context.Context) error
}
