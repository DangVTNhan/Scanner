package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionIndexes defines the indexes required for a collection
type CollectionIndexes struct {
	CollectionName string
	Indexes        []mongo.IndexModel
}

// IndexDefinitions contains all the indexes for all collections
var IndexDefinitions = []CollectionIndexes{
	{
		CollectionName: "reports",
		Indexes: []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "timestamp", Value: -1},
					{Key: "_id", Value: -1},
				},
				Options: options.Index().SetName("timestamp_desc"),
			},
		},
	},
	{
		CollectionName: "weather_cache",
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "timestamp", Value: -1}},
				Options: options.Index().SetName("timestamp_desc"),
			},
		},
	},
}

// EnsureIndexes checks and creates all required indexes for all collections
func EnsureIndexes(db *mongo.Database) error {
	log.Println("Checking and creating indexes...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, collectionIndexes := range IndexDefinitions {
		collection := db.Collection(collectionIndexes.CollectionName)

		// Check if collection exists
		if err := checkCollectionExists(ctx, db, collectionIndexes.CollectionName); err != nil {
			log.Printf("Collection %s does not exist yet, will be created when first document is inserted", collectionIndexes.CollectionName)
			continue
		}

		// Get existing indexes
		existingIndexes, err := getExistingIndexes(ctx, collection)
		if err != nil {
			return fmt.Errorf("failed to get existing indexes for collection %s: %w", collectionIndexes.CollectionName, err)
		}

		// Create missing indexes
		for _, indexModel := range collectionIndexes.Indexes {
			indexName := ""
			if indexModel.Options != nil && indexModel.Options.Name != nil {
				indexName = *indexModel.Options.Name
			}

			if !indexExists(existingIndexes, indexName) {
				log.Printf("Creating index %s for collection %s", indexName, collectionIndexes.CollectionName)

				_, err := collection.Indexes().CreateOne(ctx, indexModel)
				if err != nil {
					return fmt.Errorf("failed to create index %s for collection %s: %w", indexName, collectionIndexes.CollectionName, err)
				}

				log.Printf("Successfully created index %s for collection %s", indexName, collectionIndexes.CollectionName)
			} else {
				log.Printf("Index %s already exists for collection %s", indexName, collectionIndexes.CollectionName)
			}
		}
	}

	log.Println("All indexes have been checked and created if needed")
	return nil
}

// checkCollectionExists checks if a collection exists in the database
func checkCollectionExists(ctx context.Context, db *mongo.Database, collectionName string) error {
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		return fmt.Errorf("collection %s does not exist", collectionName)
	}

	return nil
}

// getExistingIndexes gets all existing indexes for a collection
func getExistingIndexes(ctx context.Context, collection *mongo.Collection) ([]string, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		return nil, err
	}

	var indexNames []string
	for _, index := range indexes {
		if name, ok := index["name"].(string); ok {
			indexNames = append(indexNames, name)
		}
	}

	return indexNames, nil
}

// indexExists checks if an index exists in the list of existing indexes
func indexExists(existingIndexes []string, indexName string) bool {
	for _, name := range existingIndexes {
		if name == indexName {
			return true
		}
	}
	return false
}
