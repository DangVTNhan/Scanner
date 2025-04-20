package configs

import (
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	MongoURI          string
	DatabaseName      string
	OpenWeatherAPIKey string
	Port              string
	CORS              CORSConfig
}

// CORSConfig holds the CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Default CORS allowed origins
	defaultAllowedOrigins := []string{
		"http://localhost:3000",
		"http://frontend:3000",
		"http://host.docker.internal:3000",
		"*", // Allow any origin for development
	}

	// Get allowed origins from environment variable
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowedOrigins := defaultAllowedOrigins
	if allowedOriginsEnv != "" {
		// Split by comma and trim spaces
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	}

	// CORS configuration
	corsConfig := CORSConfig{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		MaxAge:         3600,
	}

	return &Config{
		MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:      getEnv("DB_NAME", "weather_reports"),
		OpenWeatherAPIKey: os.Getenv("OPENWEATHER_API_KEY"),
		Port:              getEnv("PORT", "8080"),
		CORS:              corsConfig,
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
