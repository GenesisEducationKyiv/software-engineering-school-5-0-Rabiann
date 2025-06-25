package services

import (
	"context"

	"github.com/Rabiann/weather-mailer/internal/external"
	"github.com/Rabiann/weather-mailer/internal/models"
)

type WeatherService struct {
	weatherProvider *external.WeatherProvider
}

func NewWeatherService(weatherProvider *external.WeatherProvider) *WeatherService {
	return &WeatherService{weatherProvider}
}

func (w *WeatherService) GetWeather(city string, ctx_ context.Context) (models.Weather, error) {
	return w.weatherProvider.GetWeather(city, ctx_)

}
