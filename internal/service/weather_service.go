package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/victorgiudicissi/viacep-cloud-run/internal/model"
)

type WeatherService struct {
	apiKey     string
	baseURL    string
	httpClient HTTPClient
}

func NewWeatherService(apiKey string, client HTTPClient) *WeatherService {
	return &WeatherService{
		apiKey:  apiKey,
		baseURL: "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		httpClient: client,
	}
}

func (s *WeatherService) GetTemperature(city string) (*model.WeatherResponse, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("weather API key not configured")
	}
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf(s.baseURL, s.apiKey, encodedCity)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding weather response: %w", err)
	}

	tempC := result.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return &model.WeatherResponse{
		TempC: round(tempC, 2),
		TempF: round(tempF, 2),
		TempK: round(tempK, 2),
	}, nil
}

func round(value float64, places int) float64 {
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.*f", places, value), 64)
	return f
}
