package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/DangVTNhan/Scanner/be/configs"
)

// CORSMiddleware creates a middleware that handles Cross-Origin Resource Sharing (CORS) for the API
func CORSMiddleware(next http.Handler) http.Handler {
	// Get the application configuration
	config := configs.LoadConfig()
	corsConfig := config.CORS

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the origin from the request
		origin := r.Header.Get("Origin")

		// Check if the origin is allowed
		allowedOrigin := "*"
		if origin != "" {
			for _, allowed := range corsConfig.AllowedOrigins {
				if allowed == origin || allowed == "*" {
					allowedOrigin = origin
					break
				}
			}
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
