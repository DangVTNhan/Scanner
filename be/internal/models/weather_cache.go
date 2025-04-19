package models

import (
	"time"

	"github.com/DangVTNhan/Scanner/be/pkg/openweather"
)

// WeatherCache represents a cached weather data entry
type WeatherCache struct {
	ID          string                  `json:"id" bson:"_id,omitempty"`
	Timestamp   time.Time               `json:"timestamp" bson:"timestamp"`
	WeatherData openweather.WeatherData `json:"weatherData" bson:"weatherData"`
	CreatedAt   time.Time               `json:"createdAt" bson:"createdAt"`
}
