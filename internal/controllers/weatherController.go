package controllers

import (
	"context"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	WeatherController struct {
		weatherService WeatherService
	}

	WeatherService interface {
		GetWeather(string, context.Context) (models.Weather, error)
	}
)

func NewWeatherController(weatherService WeatherService) WeatherController {
	return WeatherController{weatherService: weatherService}
}

func (w WeatherController) GetWeather(ctx *gin.Context) {
	city, ok := ctx.GetQuery("city")
	if !ok {
		ctx.JSON(401, nil)
		return
	}

	weather, err := w.weatherService.GetWeather(city, ctx)
	if err != nil {
		ctx.JSON(400, nil)
	}

	ctx.JSON(http.StatusOK, weather)
}
