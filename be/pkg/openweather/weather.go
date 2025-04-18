package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.openweathermap.org/data/2.5"
	// Changi Airport coordinates
	latitude  = 1.3586
	longitude = 103.9899
)

// WeatherService handles interactions with the OpenWeather API
type WeatherService struct {
	apiKey string
	client *http.Client
}

// NewWeatherService creates a new instance of WeatherService
func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WeatherData represents the relevant weather data from OpenWeather API
type WeatherData struct {
	Temperature float64 // in Celsius
	Pressure    float64 // in hPa
	Humidity    float64 // in %
	CloudCover  float64 // in %
}

// apiResponse represents the response from OpenWeather API
type apiResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	Clouds struct {
		All float64 `json:"all"`
	} `json:"clouds"`
}

// GetCurrentWeather fetches the current weather for Changi Airport
func (s *WeatherService) GetCurrentWeather() (*WeatherData, error) {
	url := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=metric", baseURL, latitude, longitude, s.apiKey)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenWeather API returned non-OK status: %d", resp.StatusCode)
	}
	
	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}
	
	return &WeatherData{
		Temperature: apiResp.Main.Temp,
		Pressure:    apiResp.Main.Pressure,
		Humidity:    apiResp.Main.Humidity,
		CloudCover:  apiResp.Clouds.All,
	}, nil
}

// GetHistoricalWeather fetches historical weather data for Changi Airport
// Note: This requires a paid OpenWeather API subscription
// For a free alternative, we could store our own historical data
func (s *WeatherService) GetHistoricalWeather(timestamp time.Time) (*WeatherData, error) {
	// For demonstration purposes, we'll return mock data
	// In a real application, you would use the OpenWeather historical data API
	// or implement a solution to store and retrieve your own historical data
	
	// This is a simplified implementation
	return &WeatherData{
		Temperature: 30.0,
		Pressure:    1013.0,
		Humidity:    70.0,
		CloudCover:  40.0,
	}, nil
}
