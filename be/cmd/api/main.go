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
	"github.com/DangVTNhan/Scanner/be/internal/handlers"
	"github.com/DangVTNhan/Scanner/be/internal/middleware"
	"github.com/DangVTNhan/Scanner/be/internal/services"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	if config.OpenWeatherAPIKey == "" {
		log.Fatal("OPENWEATHER_API_KEY environment variable is required")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Initialize services
	db := client.Database(config.DatabaseName)
	weatherService := openweather.NewWeatherService(config.OpenWeatherAPIKey)
	reportService := services.NewReportService(db, weatherService)

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
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	srv.Shutdown(ctx)

	fmt.Println("Server gracefully stopped")
}
