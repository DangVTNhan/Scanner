package mongodb

import (
	"context"
	"errors"
	"github.com/DangVTNhan/Scanner/be/internal/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCursor is a mock implementation of ICursor
type MockCursor struct {
	mock.Mock
	results []models.WeatherReport
	current int
}

func NewMockCursor(results []models.WeatherReport) *MockCursor {
	return &MockCursor{
		results: results,
		current: 0,
	}
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

func (m *MockCursor) All(ctx context.Context, results interface{}) error {
	args := m.Called(ctx, results)

	// If we have results to return and a pointer to a slice is provided
	if len(m.results) > 0 {
		if resultsPtr, ok := results.(*[]models.WeatherReport); ok {
			*resultsPtr = m.results
		}
	}

	return args.Error(0)
}

// MockCollection is a mock implementation of ICollection
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) ISingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(ISingleResult)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (ICursor, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(ICursor), args.Error(1)
}

func (m *MockCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

// MockSingleResult is a mock implementation of ISingleResult
type MockSingleResult struct {
	mock.Mock
	err error
	doc interface{}
}

func (m *MockSingleResult) Decode(v interface{}) error {
	if m.err != nil {
		return m.err
	}
	if m.doc == nil {
		return mongo.ErrNoDocuments
	}

	// If doc is provided, copy its values to v
	switch doc := m.doc.(type) {
	case *models.WeatherReport:
		if report, ok := v.(*models.WeatherReport); ok {
			*report = *doc
			return nil
		}
	case *models.WeatherCache:
		if cache, ok := v.(*models.WeatherCache); ok {
			*cache = *doc
			return nil
		}
	}
	return errors.New("could not decode value")
}

func (m *MockSingleResult) Err() error {
	return m.err
}

// MockDatabase is a mock implementation of IDatabase
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Collection(name string, opts ...*options.CollectionOptions) ICollection {
	args := m.Called(name, opts)
	return args.Get(0).(ICollection)
}

func NewMockSingleResult(err error, doc interface{}) *MockSingleResult {
	return &MockSingleResult{
		err: err,
		doc: doc,
	}
}
