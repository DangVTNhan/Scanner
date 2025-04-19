package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoWeatherCacheRepository implements the IWeatherCacheRepository interface for MongoDB
type MongoWeatherCacheRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoWeatherCacheRepository creates a new instance of MongoWeatherCacheRepository
func NewMongoWeatherCacheRepository(db *mongo.Database) repository.IWeatherCacheRepository {
	return &MongoWeatherCacheRepository{
		db:         db,
		collection: db.Collection("weather_cache"),
	}
}

// SaveWeatherCache saves a weather data cache entry
func (r *MongoWeatherCacheRepository) SaveWeatherCache(ctx context.Context, cache *models.WeatherCache) (string, error) {
	result, err := r.collection.InsertOne(ctx, cache)
	if err != nil {
		return "", fmt.Errorf("failed to save weather cache: %w", err)
	}

	// Convert ObjectID to string
	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID.Hex(), nil
}

// FindLatestWeatherCache retrieves the latest valid weather cache entry
func (r *MongoWeatherCacheRepository) FindLatestWeatherCache(ctx context.Context) (*models.WeatherCache, error) {
	now := time.Now()

	// Find the latest non-expired cache entry
	filter := bson.M{
		"expiresAt": bson.M{"$gt": now},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	var cache models.WeatherCache
	err := r.collection.FindOne(ctx, filter, opts).Decode(&cache)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No cache found, not an error
		}
		return nil, fmt.Errorf("failed to retrieve weather cache: %w", err)
	}

	return &cache, nil
}

// FindWeatherCacheByTimestamp retrieves a weather cache entry by timestamp within a time window
func (r *MongoWeatherCacheRepository) FindWeatherCacheByTimestamp(ctx context.Context, timestamp time.Time, windowMinutes ...int) (*models.WeatherCache, error) {
	window := 10 // Default window in minutes
	if len(windowMinutes) > 0 {
		window = windowMinutes[0]
	}
	// Calculate the time window
	windowStart := timestamp.Add(-time.Duration(window) * time.Minute)
	windowEnd := timestamp.Add(time.Duration(window) * time.Minute)

	// Find a cache entry within the time window that hasn't expired
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": windowStart,
			"$lte": windowEnd,
		},
	}

	var cache models.WeatherCache
	err := r.collection.FindOne(ctx, filter).Decode(&cache)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No cache found, not an error
		}
		return nil, fmt.Errorf("failed to retrieve weather cache: %w", err)
	}

	return &cache, nil
}

// DeleteExpiredCaches removes expired cache entries
func (r *MongoWeatherCacheRepository) DeleteExpiredCaches(ctx context.Context) error {
	now := time.Now()

	filter := bson.M{
		"expiresAt": bson.M{"$lte": now},
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete expired caches: %w", err)
	}

	return nil
}
