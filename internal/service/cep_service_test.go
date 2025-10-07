package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewCEPService(t *testing.T) {
	t.Parallel()
	service := NewCEPService(nil)
	assert.NotNil(t, service)
	assert.Equal(t, "https://viacep.com.br/ws/%s/json/", service.baseURL)
	assert.NotNil(t, service.httpClient)
}

func TestGetCity_ValidCEP(t *testing.T) {
	t.Parallel()
	mockClient := new(MockHTTPClient)
	
	expectedURL := "https://viacep.com.br/ws/01001000/json/"
	_, _ = http.NewRequest("GET", expectedURL, nil)
	
	responseBody := `{"cep": "01001-000", "localidade": "São Paulo"}`
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
	
	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL.String() == expectedURL
	})).Return(response, nil)

	service := NewCEPService(mockClient)
	city, err := service.GetCity("01001000")

	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)
	mockClient.AssertExpectations(t)
}

func TestGetCity_InvalidCEPFormat(t *testing.T) {
	t.Parallel()
	mockClient := new(MockHTTPClient)
	service := NewCEPService(mockClient)

	_, err := service.GetCity("12345")
	assert.EqualError(t, err, "invalid zipcode")

	_, err = service.GetCity("abcde123")
	assert.EqualError(t, err, "invalid zipcode")

	mockClient.AssertNotCalled(t, "Do", mock.Anything)
}

func TestGetCity_CEPNotFound(t *testing.T) {
	t.Parallel()
	mockClient := new(MockHTTPClient)
	
	expectedURL := "https://viacep.com.br/ws/99999999/json/"
	_, _ = http.NewRequest("GET", expectedURL, nil)
	
	responseBody := `{"erro": "true"}`
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}
	
	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL.String() == expectedURL
	})).Return(response, nil)

	service := NewCEPService(mockClient)
	_, err := service.GetCity("99999999")
	assert.EqualError(t, err, "can not find zipcode")
	mockClient.AssertExpectations(t)
}

func TestGetCity_InvalidJSONResponse(t *testing.T) {
	t.Parallel()
	mockClient := new(MockHTTPClient)
	
	expectedURL := "https://viacep.com.br/ws/01001000/json/"
	_, _ = http.NewRequest("GET", expectedURL, nil)
	
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("invalid json")),
	}
	
	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL.String() == expectedURL
	})).Return(response, nil)

	service := NewCEPService(mockClient)
	_, err := service.GetCity("01001000")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding CEP response")
	mockClient.AssertExpectations(t)
}

func TestGetCity_HTTPError(t *testing.T) {
	t.Parallel()
	mockClient := new(MockHTTPClient)
	
	expectedURL := "https://viacep.com.br/ws/01001000/json/"
	_, _ = http.NewRequest("GET", expectedURL, nil)
	
	expectedErr := errors.New("connection error")
	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL.String() == expectedURL
	})).Return((*http.Response)(nil), expectedErr)

	service := NewCEPService(mockClient)
	_, err := service.GetCity("01001000")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching CEP")
	mockClient.AssertExpectations(t)
}

func TestGetCity_InvalidStatusCode(t *testing.T) {
	t.Parallel()
	mockClient := new(MockHTTPClient)
	
	expectedURL := "https://viacep.com.br/ws/01001000/json/"
	_, _ = http.NewRequest("GET", expectedURL, nil)
	
	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	
	mockClient.On("Do", mock.MatchedBy(func(r *http.Request) bool {
		return r.URL.String() == expectedURL
	})).Return(response, nil)

	service := NewCEPService(mockClient)
	_, err := service.GetCity("01001000")

	assert.EqualError(t, err, "invalid zipcode")
	mockClient.AssertExpectations(t)
}
