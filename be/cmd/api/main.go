package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DangVTNhan/Scanner/be/configs"
	_ "github.com/DangVTNhan/Scanner/be/docs" // Import swagger docs
	"github.com/DangVTNhan/Scanner/be/internal/database"
	"github.com/DangVTNhan/Scanner/be/internal/handlers"
	"github.com/DangVTNhan/Scanner/be/internal/middleware"
	"github.com/DangVTNhan/Scanner/be/internal/models/repository/mongodb"
	"github.com/DangVTNhan/Scanner/be/internal/services"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Changi Airport Weather Report API
// @version         1.0
// @description     API for generating and retrieving weather reports for Changi Airport
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  nhan.dangviettrung@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api
func main() {
	// Load configuration
	config := configs.LoadConfig()

	if config.OpenWeatherAPIKey == "" {
		log.Fatal("OPENWEATHER_API_KEY environment variable is required")
	}

	// Connect to MongoDB and initialize database with indexes
	client, db, err := database.InitDatabase(config.MongoURI, config.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Wrap the database with our interface wrapper
	dbWrapper := mongodb.NewMongoDatabaseWrapper(db)

	// Initialize repositories
	reportRepository := mongodb.NewMongoReportRepository(dbWrapper)
	weatherCacheRepository := mongodb.NewMongoWeatherCacheRepository(dbWrapper)

	// Initialize weather service with caching
	weatherService := openweather.NewWeatherService(config.OpenWeatherAPIKey)

	// Initialize services with repositories
	reportService := services.NewReportService(reportRepository, weatherCacheRepository, weatherService)

	// Initialize handlers
	reportHandler := handlers.NewReportHandler(reportService)

	// Set up router
	router := mux.NewRouter()

	// Apply CORS middleware - must be added before routes
	router.Use(middleware.CORSMiddleware)

	// API routes
	router.HandleFunc("/api/reports", reportHandler.GenerateReport).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/reports", reportHandler.GetAllReports).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/reports/paginated", reportHandler.GetPaginatedReports).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/reports/{id}", reportHandler.GetReportByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/reports/compare", reportHandler.CompareReports).Methods("POST", "OPTIONS")

	// Swagger documentation - only available in dev/stg environments
	if config.IsSwaggerEnabled() {
		fmt.Println("Swagger UI enabled at /swagger/index.html")
		router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		))
	}

	// Start server
	addr := ":" + config.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in a goroutine so that it doesn't block
	go func() {
		fmt.Printf("Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	srv.Shutdown(ctx)

	fmt.Println("Server gracefully stopped")
}
