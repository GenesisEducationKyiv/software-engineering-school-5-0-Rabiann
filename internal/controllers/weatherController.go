package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/gin-gonic/gin"
)

type (
	WeatherController struct {
		weatherService WeatherService
	}

	WeatherService interface {
		GetWeather(string, context.Context, context.CancelFunc) (models.Weather, error)
	}
)

func NewWeatherController(weatherService WeatherService) WeatherController {
	return WeatherController{weatherService: weatherService}
}

func (w WeatherController) GetWeather(ctx *gin.Context) {
	ctx_, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	city, ok := ctx.GetQuery("city")
	if !ok {
		ctx.JSON(400, nil)
		return
	}

	weather, err := w.weatherService.GetWeather(city, ctx_, cancel)
	if err != nil {
		ctx.JSON(400, nil)
	}

	ctx.JSON(http.StatusOK, weather)
}
