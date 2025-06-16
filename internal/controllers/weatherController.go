package controllers

import (
	"net/http"

	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	weatherService services.WeatherServer
}

func NewWeatherController(weatherService *services.WeatherService) WeatherController {
	return WeatherController{weatherService: weatherService}
}

func (w WeatherController) GetWeather(ctx *gin.Context) {
	city, ok := ctx.GetQuery("city")
	if !ok {
		ctx.JSON(400, nil)
		return
	}

	weather, err := w.weatherService.GetWeather(city)
	if err != nil {
		ctx.JSON(400, nil)
	}

	ctx.JSON(http.StatusOK, weather)
}
