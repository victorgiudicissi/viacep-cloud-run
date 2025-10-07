package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/viacep-cloud-run/internal/model"
	"github.com/victorgiudicissi/viacep-cloud-run/internal/service"
)

type Handler struct {
	weatherService *service.WeatherService
	cepService    *service.CEPService
}

func NewHandler(weatherService *service.WeatherService, cepService *service.CEPService) *Handler {
	return &Handler{
		weatherService: weatherService,
		cepService:    cepService,
	}
}

func (h *Handler) GetWeatherByCEP(c *gin.Context) {
	cep := c.Param("cep")

	city, err := h.cepService.GetCity(cep)
	if err != nil {
		switch err.Error() {
		case "invalid zipcode":
			c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "invalid zipcode"})
		case "can not find zipcode":
			c.JSON(http.StatusNotFound, model.ErrorResponse{Message: "can not find zipcode"})
		default:
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "internal server error"})
		}
		return
	}

	weather, err := h.weatherService.GetTemperature(city)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "error getting weather data"})
		return
	}

	c.JSON(http.StatusOK, weather)
}
