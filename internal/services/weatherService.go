package services

import (
	"context"
	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	WeatherService struct {
		weatherProvider WeatherProvider
	}

	WeatherProvider interface {
		GetWeather(city string, ctx context.Context) (models.Weather, error)
	}
)

func NewWeatherService(weatherProvider WeatherProvider) *WeatherService {
	return &WeatherService{weatherProvider}
}

func (w *WeatherService) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	return w.weatherProvider.GetWeather(city, ctx)

}
