# CEP Weather Service

A Go service that retrieves weather information based on Brazilian CEP (ZIP code).

## Features

- Get current weather by CEP (Brazilian ZIP code)
- Returns temperature in Celsius, Fahrenheit, and Kelvin
- Containerized with Docker

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- A WeatherAPI.com API key

## Getting Started

1. Clone the repository
2. Copy your Weather API Key to `.env`:
3. Build and run with Docker Compose:
   ```bash
   make run
   ```

## API Endpoints

### Get Weather by CEP

```
GET /weather/{cep}
```

#### Successful Response (200 OK)
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

#### Error Responses

- **404 Not Found**: CEP not found
- **422 Unprocessable Entity**: Invalid CEP format
- **500 Internal Server Error**: Server error

## Development

### Running Tests

```bash
make test
```

### Running Locally

1. Set up environment variables:
   ```bash
   export WEATHER_API_KEY=your_api_key_here
   export PORT=8080
   ```

2. Run the application:
   ```bash
   make run
   ```

3. Call a valid CEP:
```
curl --location 'http://localhost:8080/weather/13574560'
```

The response will be something like:
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

### With Google Cloud Run:
```
curl --location 'https://viacep-cloud-run-83741703367.us-central1.run.app/weather/13574560'
```

The response will be something like:
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```
