package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CEPService struct {
	baseURL    string
	httpClient HTTPClient
}

func NewCEPService(client HTTPClient) *CEPService {
	if client == nil {
		client = &http.Client{}
	}
	return &CEPService{
		baseURL:    "https://viacep.com.br/ws/%s/json/",
		httpClient: client,
	}
}

func (s *CEPService) GetCity(cep string) (string, error) {
	re := regexp.MustCompile(`^\d{8}$`)
	if !re.MatchString(cep) {
		return "", errors.New("invalid zipcode")
	}

	url := fmt.Sprintf(s.baseURL, cep)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching CEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", errors.New("invalid zipcode")
	}

	var result struct {
		Localidade string `json:"localidade"`
		Erro       string   `json:"erro"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding CEP response: %w", err)
	}

	if result.Erro == "true" || result.Localidade == "" {
		return "", errors.New("can not find zipcode")
	}

	return result.Localidade, nil
}
