package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHandler is a simple handler that records if it was called
type mockHandler struct {
	called bool
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.called = true
	w.WriteHeader(http.StatusOK)
}

// TestCORSMiddleware_OptionsRequest tests that OPTIONS requests are handled correctly
func TestCORSMiddleware_OptionsRequest(t *testing.T) {
	// Create a mock handler
	mock := &mockHandler{}

	// Save original environment variable
	originalAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Set environment variable for test
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,*")

	// Restore original environment variable when test completes
	defer func() { os.Setenv("CORS_ALLOWED_ORIGINS", originalAllowedOrigins) }()

	// Create a request with OPTIONS method
	req := httptest.NewRequest("OPTIONS", "/api/reports", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create the middleware
	middleware := CORSMiddleware(mock)

	// Call the middleware
	middleware.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the mock handler was not called (OPTIONS requests should be handled by the middleware)
	assert.False(t, mock.called)

	assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

// TestCORSMiddleware_NonOptionsRequest tests that non-OPTIONS requests are passed to the next handler
func TestCORSMiddleware_NonOptionsRequest(t *testing.T) {
	// Create a mock handler
	mock := &mockHandler{}

	// Save original environment variable
	originalAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Set environment variable for test
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,*")

	// Restore original environment variable when test completes
	defer func() { os.Setenv("CORS_ALLOWED_ORIGINS", originalAllowedOrigins) }()

	// Create a request with GET method
	req := httptest.NewRequest("GET", "/api/reports", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create the middleware
	middleware := CORSMiddleware(mock)

	// Call the middleware
	middleware.ServeHTTP(rr, req)

	// Check that the mock handler was called
	assert.True(t, mock.called)

	assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

// TestCORSMiddleware_AllowedOrigin tests that allowed origins are handled correctly
func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	// Create a mock handler
	mock := &mockHandler{}

	// Save original environment variable
	originalAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Set environment variable for test with specific allowed origins
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://frontend:3000,*")

	// Restore original environment variable when test completes
	defer func() { os.Setenv("CORS_ALLOWED_ORIGINS", originalAllowedOrigins) }()

	// Test cases for different origins
	testCases := []struct {
		name           string
		origin         string
		expectedOrigin string
	}{
		{
			name:           "Allowed specific origin",
			origin:         "http://localhost:3000",
			expectedOrigin: "http://localhost:3000",
		},
		{
			name:           "Allowed frontend origin",
			origin:         "http://frontend:3000",
			expectedOrigin: "http://frontend:3000",
		},
		{
			name:           "Disallowed origin",
			origin:         "https://www.stengg.com/",
			expectedOrigin: "https://www.stengg.com/",
		},
		{
			name:           "No origin",
			origin:         "",
			expectedOrigin: "*", // Falls back to wildcard
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the mock
			mock.called = false

			// Create a request
			req := httptest.NewRequest("GET", "/api/reports", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create the middleware
			middleware := CORSMiddleware(mock)

			// Call the middleware
			middleware.ServeHTTP(rr, req)

			// Check that the mock handler was called
			assert.True(t, mock.called)

			// Check the Access-Control-Allow-Origin header
			assert.Equal(t, tc.expectedOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

// TestCORSMiddleware_HeadersSet tests that all required CORS headers are set
func TestCORSMiddleware_HeadersSet(t *testing.T) {
	// Create a mock handler
	mock := &mockHandler{}

	// Save original environment variable
	originalAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Set environment variable for test
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,*")

	// Restore original environment variable when test completes
	defer func() { os.Setenv("CORS_ALLOWED_ORIGINS", originalAllowedOrigins) }()

	// Create a request
	req := httptest.NewRequest("GET", "/api/reports", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create the middleware
	middleware := CORSMiddleware(mock)

	// Call the middleware
	middleware.ServeHTTP(rr, req)

	// Check all required CORS headers
	assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization, X-Requested-With", rr.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "3600", rr.Header().Get("Access-Control-Max-Age"))
}
