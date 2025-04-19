package mongodb

import (
	"context"
	"fmt"

	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/DangVTNhan/Scanner/be/internal/models/repository"
	"github.com/DangVTNhan/Scanner/be/internal/models/request"
	"github.com/DangVTNhan/Scanner/be/internal/models/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoReportRepository implements the IReportRepository interface for MongoDB
type MongoReportRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoReportRepository creates a new instance of MongoReportRepository
func NewMongoReportRepository(db *mongo.Database) repository.IReportRepository {
	return &MongoReportRepository{
		db:         db,
		collection: db.Collection("reports"),
	}
}

// InsertReport inserts a new weather report into the database
func (r *MongoReportRepository) InsertReport(ctx context.Context, report *models.WeatherReport) (string, error) {
	result, err := r.collection.InsertOne(ctx, report)
	if err != nil {
		return "", fmt.Errorf("failed to save report: %w", err)
	}

	// Convert ObjectID to string
	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID.Hex(), nil
}

// FindAllReports retrieves all weather reports
func (r *MongoReportRepository) FindAllReports(ctx context.Context) ([]models.WeatherReport, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
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

// FindPaginatedReports retrieves weather reports with pagination and filtering
func (r *MongoReportRepository) FindPaginatedReports(ctx context.Context, req *request.PaginatedReportsRequest) (*response.PaginatedReportsResponse, error) {
	// Set default limit if not provided
	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}

	// Build the filter
	filter := bson.M{}

	// Add time range filter if provided
	if !req.FromTime.IsZero() || !req.ToTime.IsZero() {
		timeFilter := bson.M{}
		if !req.FromTime.IsZero() {
			timeFilter["$gte"] = req.FromTime
		}
		if !req.ToTime.IsZero() {
			timeFilter["$lte"] = req.ToTime
		}
		filter["timestamp"] = timeFilter
	}

	// Add cursor-based pagination if lastID is provided
	if req.LastID != "" {
		// Try to convert to ObjectID first
		lastObjectID, err := primitive.ObjectIDFromHex(req.LastID)
		if err != nil {
			// If not a valid ObjectID, use string ID directly
			filter["_id"] = bson.M{"$lt": req.LastID}
		} else {
			// If it's a valid ObjectID, use it
			filter["_id"] = bson.M{"$lt": lastObjectID}
		}
	}

	// Set up options for sorting and limiting
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(limit + 1)) // Fetch one extra to check if there are more

	// Execute the query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve reports: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results
	var reports []models.WeatherReport
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, fmt.Errorf("failed to decode reports: %w", err)
	}

	// Check if there are more results
	hasMore := false
	if len(reports) > limit {
		hasMore = true
		reports = reports[:limit] // Remove the extra item
	}

	// Calculate pagination metadata
	response := &response.PaginatedReportsResponse{
		Reports:     reports,
		HasMore:     hasMore,
		CurrentPage: 1, // Default to 1 for cursor-based pagination
	}

	// If we're using cursor-based pagination, calculate the current page
	if req.LastID != "" {
		response.CurrentPage = 2 // At least page 2 if lastID is provided
	}

	// Calculate from and to numbers
	response.FromNumber = (response.CurrentPage-1)*limit + 1
	response.ToNumber = response.FromNumber + len(reports) - 1

	// Get total count only if not filtered
	if !req.IsFiltered {
		totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("failed to count reports: %w", err)
		}
		response.TotalCount = int(totalCount)
	}

	return response, nil
}

// FindReportByID retrieves a weather report by its ID
func (r *MongoReportRepository) FindReportByID(ctx context.Context, id string) (*models.WeatherReport, error) {
	// Check if the ID is a valid ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// If not a valid ObjectID, try to find by string ID directly
		var report models.WeatherReport
		err = r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&report)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("report not found")
			}
			return nil, fmt.Errorf("failed to retrieve report: %w", err)
		}
		return &report, nil
	}

	// If it's a valid ObjectID, search by ObjectID
	var report models.WeatherReport
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("report not found")
		}
		return nil, fmt.Errorf("failed to retrieve report: %w", err)
	}

	return &report, nil
}

// CountReports counts the total number of reports
func (r *MongoReportRepository) CountReports(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count reports: %w", err)
	}
	return count, nil
}
