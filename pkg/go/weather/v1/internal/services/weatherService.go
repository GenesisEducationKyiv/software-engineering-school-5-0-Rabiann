package services

import (
	"context"

	"github.com/Rabiann/weather-mailer/pkg/go/weather/internal/models"
)

type (
	WeatherService struct {
		weatherProvider WeatherProvider
		cache           CacheServicer
	}

	WeatherProvider interface {
		GetWeather(city string, ctx context.Context) (models.Weather, error)
	}

	CacheServicer interface {
		Read(string) (models.Weather, bool)
		Write(string, models.Weather) error
	}
)

func NewWeatherService(weatherProvider WeatherProvider, cache CacheServicer) *WeatherService {
	return &WeatherService{weatherProvider, cache}
}

func (w *WeatherService) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	var weather models.Weather

	cache := w.cache

	weather, ok := cache.Read(city)

	if !ok {
		weather, err := w.weatherProvider.GetWeather(city, ctx)
		if err != nil {
			return weather, err
		}

		if err = cache.Write(city, weather); err != nil {
			return weather, err
		}

		return weather, nil
	}

	return weather, nil
}
