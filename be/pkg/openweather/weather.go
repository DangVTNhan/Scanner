package openweather

import (
	"encoding/json"
	"fmt"
	"github.com/DangVTNhan/Scanner/be/pkg/openweather/response"
	"net/http"
	"time"

	"golang.org/x/sync/singleflight"
)

type IWeatherService interface {
	// GetCurrentWeather fetches the current weather
	GetCurrentWeather() (*WeatherData, error)

	// GetHistoricalWeather fetches historical weather data
	GetHistoricalWeather(timestamp time.Time) (*WeatherData, error)
}

type WeatherService struct {
	apiKey string
	client *http.Client
	sfg    singleflight.Group
}


const (
	baseURL = "https://api.openweathermap.org/data/3.0/onecall"
	// Changi Airport coordinates
	latitude  = 1.3586
	longitude = 103.9899
)

// WeatherService handles interactions with the OpenWeather API

// NewWeatherService creates a new instance of WeatherService
func NewWeatherService(apiKey string) IWeatherService {
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
// GetCurrentWeather fetches the current weather for Changi Airport
// Uses singleflight to deduplicate concurrent requests
func (s *WeatherService) GetCurrentWeather() (*WeatherData, error) {
	// Use a constant key for current weather since it's the same for all requests
	const sfKey = "current_weather"

	// Use singleflight to deduplicate concurrent requests
	result, err, _ := s.sfg.Do(sfKey, func() (interface{}, error) {
		url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", baseURL, latitude, longitude, s.apiKey)

		resp, err := s.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch weather data: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("OpenWeather API returned non-OK status: %d", resp.StatusCode)
		}

		var apiResp response.GetCurrentWeatherResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("failed to decode API response: %w", err)
		}

		return &WeatherData{
			Temperature: apiResp.Current.Temp,
			Pressure:    apiResp.Current.Pressure,
			Humidity:    apiResp.Current.Humidity,
			CloudCover:  apiResp.Current.Clouds,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	// Type assertion to convert the interface{} result back to *WeatherData
	return result.(*WeatherData), nil
}

// GetHistoricalWeather fetches historical weather data for Changi Airport
// Note: This requires a paid OpenWeather API subscription
// For a free alternative, we could store our own historical data
// Uses singleflight to deduplicate concurrent requests for the same timestamp
func (s *WeatherService) GetHistoricalWeather(timestamp time.Time) (*WeatherData, error) {
	// Use timestamp as the key for singleflight to deduplicate concurrent requests for the same timestamp
	sfKey := fmt.Sprintf("historical_weather_%d", timestamp.Unix())

	// Use singleflight to deduplicate concurrent requests
	result, err, _ := s.sfg.Do(sfKey, func() (interface{}, error) {
		url := fmt.Sprintf("%s/timemachine?lat=%f&lon=%f&dt=%d&appid=%s&units=metric", baseURL, latitude, longitude, timestamp.Unix(), s.apiKey)

		resp, err := s.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch weather data: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("OpenWeather API returned non-OK status: %d", resp.StatusCode)
		}

		var apiResp response.GetHistoricalTimeResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("failed to decode API response: %w", err)
		}

		if len(apiResp.Data) == 0 {
			return nil, fmt.Errorf("no historical data found for the given timestamp")
		}
		// Use the first data point in the response
		data := apiResp.Data[0]
		return &WeatherData{
			Temperature: data.Temp,
			Pressure:    data.Pressure,
			Humidity:    data.Humidity,
			CloudCover:  data.Clouds,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	// Type assertion to convert the interface{} result back to *WeatherData
	return result.(*WeatherData), nil
}
