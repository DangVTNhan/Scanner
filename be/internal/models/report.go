package models

import (
	"time"
)

type WeatherReport struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Temperature float64   `json:"temperature" bson:"temperature"` // in Celsius
	Pressure    float64   `json:"pressure" bson:"pressure"`       // in hPa
	Humidity    float64   `json:"humidity" bson:"humidity"`       // in %
	CloudCover  float64   `json:"cloudCover" bson:"cloudCover"`   // in %
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}
