package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to MongoDB and returns a database instance
func Connect(ctx context.Context, mongoURI, databaseName string) (*mongo.Database, error) {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	
	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	
	log.Println("Connected to MongoDB successfully")
	
	// Get database instance
	db := client.Database(databaseName)
	
	// Ensure indexes
	if err := EnsureIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to ensure indexes: %w", err)
	}
	
	return db, nil
}

// Disconnect closes the MongoDB connection
func Disconnect(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	
	log.Println("Disconnected from MongoDB successfully")
	return nil
}

// InitDatabase initializes the database connection and ensures indexes
func InitDatabase(mongoURI, databaseName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	
	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	
	log.Println("Connected to MongoDB successfully")
	
	// Get database instance
	db := client.Database(databaseName)
	
	// Ensure indexes
	if err := EnsureIndexes(db); err != nil {
		return nil, nil, fmt.Errorf("failed to ensure indexes: %w", err)
	}
	
	return client, db, nil
}
