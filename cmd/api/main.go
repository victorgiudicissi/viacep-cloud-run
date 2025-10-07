package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/victorgiudicissi/viacep-cloud-run/internal/handler"
	"github.com/victorgiudicissi/viacep-cloud-run/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherService := service.NewWeatherService(os.Getenv("WEATHER_API_KEY"), httpClient)
	cepService := service.NewCEPService(httpClient)

	h := handler.NewHandler(weatherService, cepService)

	r := gin.Default()

	r.GET("/weather/:cep", h.GetWeatherByCEP)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
