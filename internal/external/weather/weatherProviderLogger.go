package weather

import (
	"context"
	"log"

	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	WeatherProviderLogger struct {
		provider *WeatherProvider
	}

	Logger interface {
		Info(log string)
		Error(log string)
	}
)

func NewWeatherProviderLogger(provider *WeatherProvider) *WeatherProviderLogger {
	return &WeatherProviderLogger{provider}
}

func (w *WeatherProviderLogger) GetWeather(city string, ctx_ context.Context) (models.Weather, error) {
	provider := w.provider
	resp, err := provider.GetWeather(city, ctx_)
	if err != nil {
		log.Fatalf("%s - Response: %+v \tError: %s", provider.Name(), resp, err)
		return resp, err
	}

	log.Printf("%s - Response: %+v", provider.Name(), resp)
	return resp, err
}
