package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/victorgiudicissi/viacep-cloud-run/internal/model"
)

func TestGetTemperature_Success(t *testing.T) {
	t.Parallel()
	mockClient := &MockHTTPClient{}
	
	responseBody := `{"current": {"temp_c": 25.0}}`
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}

	mockClient.On("Do", mock.Anything).Return(response, nil)

	service := &WeatherService{
		apiKey:     "test-api-key",
		baseURL:    "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		httpClient: mockClient,
	}

	result, err := service.GetTemperature("Sao Paulo")

	assert.NoError(t, err)
	expected := &model.WeatherResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298,
	}
	assert.Equal(t, expected.TempC, result.TempC)
	assert.Equal(t, expected.TempF, result.TempF)
	assert.Equal(t, expected.TempK, result.TempK)
}

func TestGetTemperature_MissingAPIKey(t *testing.T) {
	t.Parallel()
	service := &WeatherService{
		apiKey: "",
	}

	_, err := service.GetTemperature("Sao Paulo")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "weather API key not configured")
}

func TestGetTemperature_HTTPError(t *testing.T) {
	t.Parallel()
	mockClient := &MockHTTPClient{}
	mockClient.On("Do", mock.Anything).Return(
		(*http.Response)(nil),
		errors.New("http error"),
	)

	service := &WeatherService{
		apiKey:     "test-api-key",
		baseURL:    "http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		httpClient: mockClient,
	}

	_, err := service.GetTemperature("Sao Paulo")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching weather data: http error")
}

func TestRound_ToTwoDecimalPlaces(t *testing.T) {
	t.Parallel()
	result := round(3.14159, 2)
	assert.InDelta(t, 3.14, result, 0.001)
}

func TestRound_ToZeroDecimalPlaces(t *testing.T) {
	t.Parallel()
	result := round(3.6, 0)
	assert.InDelta(t, 4.0, result, 0.001)
}
