package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WeatherReport struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	Temperature float64            `json:"temperature" bson:"temperature"` // in Celsius
	Pressure    float64            `json:"pressure" bson:"pressure"`       // in hPa
	Humidity    float64            `json:"humidity" bson:"humidity"`       // in %
	CloudCover  float64            `json:"cloudCover" bson:"cloudCover"`   // in %
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}
