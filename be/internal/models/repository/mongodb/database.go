package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ICursor defines the interface for MongoDB cursor operations
type ICursor interface {
	// Close closes the cursor
	Close(ctx context.Context) error

	// Next advances the cursor to the next document
	Next(ctx context.Context) bool

	// Decode decodes the current document into the provided value
	Decode(v interface{}) error

	// All decodes all documents from the cursor into the provided slice
	All(ctx context.Context, results interface{}) error
}

// ICollection defines the interface for MongoDB collection operations
type ICollection interface {
	// InsertOne inserts a single document into the collection
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// FindOne finds a single document in the collection
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) ISingleResult

	// Find finds all documents in the collection that match the filter
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (ICursor, error)

	// DeleteMany deletes all documents from the collection that match the filter
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// CountDocuments returns the number of documents in the collection that match the filter
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

// ISingleResult defines the interface for MongoDB single result operations
type ISingleResult interface {
	// Decode decodes the result into the provided value
	Decode(v interface{}) error

	// Err returns any error encountered during the operation
	Err() error
}

// IDatabase defines the interface for database operations
type IDatabase interface {
	// Collection returns a handle to a MongoDB collection
	Collection(name string, opts ...*options.CollectionOptions) ICollection
}

// MongoCollectionWrapper wraps a mongo.Collection to implement ICollection
type MongoCollectionWrapper struct {
	coll *mongo.Collection
}

func (w *MongoCollectionWrapper) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return w.coll.InsertOne(ctx, document, opts...)
}

func (w *MongoCollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) ISingleResult {
	return &MongoSingleResultWrapper{result: w.coll.FindOne(ctx, filter, opts...)}
}

func (w *MongoCollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (ICursor, error) {
	cursor, err := w.coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &MongoCursorWrapper{cursor: cursor}, nil
}

func (w *MongoCollectionWrapper) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return w.coll.DeleteMany(ctx, filter, opts...)
}

func (w *MongoCollectionWrapper) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return w.coll.CountDocuments(ctx, filter, opts...)
}

// MongoSingleResultWrapper wraps a mongo.SingleResult to implement ISingleResult
type MongoSingleResultWrapper struct {
	result *mongo.SingleResult
}

func (w *MongoSingleResultWrapper) Decode(v interface{}) error {
	return w.result.Decode(v)
}

func (w *MongoSingleResultWrapper) Err() error {
	return w.result.Err()
}

// MongoDatabaseWrapper wraps a mongo.Database to implement IDatabase
type MongoDatabaseWrapper struct {
	db *mongo.Database
}

func NewMongoDatabaseWrapper(db *mongo.Database) IDatabase {
	return &MongoDatabaseWrapper{db: db}
}

func (w *MongoDatabaseWrapper) Collection(name string, opts ...*options.CollectionOptions) ICollection {
	return &MongoCollectionWrapper{coll: w.db.Collection(name, opts...)}
}

// MongoCursorWrapper wraps a mongo.Cursor to implement ICursor
type MongoCursorWrapper struct {
	cursor *mongo.Cursor
}

func (w *MongoCursorWrapper) Close(ctx context.Context) error {
	return w.cursor.Close(ctx)
}

func (w *MongoCursorWrapper) Next(ctx context.Context) bool {
	return w.cursor.Next(ctx)
}

func (w *MongoCursorWrapper) Decode(v interface{}) error {
	return w.cursor.Decode(v)
}

func (w *MongoCursorWrapper) All(ctx context.Context, results interface{}) error {
	return w.cursor.All(ctx, results)
}

// Ensure that our wrappers implement the interfaces
var _ ICollection = (*MongoCollectionWrapper)(nil)
var _ ISingleResult = (*MongoSingleResultWrapper)(nil)
var _ IDatabase = (*MongoDatabaseWrapper)(nil)
var _ ICursor = (*MongoCursorWrapper)(nil)
