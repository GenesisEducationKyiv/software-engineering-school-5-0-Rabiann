package weather

import (
	"context"
	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	WeatherProviderWithLaydown struct {
		curr WeatherProvider
		next LaydownWeatherProvider
	}

	LaydownWeatherProvider interface {
		GetWeather(city string, ctx context.Context) (models.Weather, error)
		SetNext(next LaydownWeatherProvider)
	}
)

func NewWeatherProviderWithLaydown() *WeatherProviderWithLaydown {
	return &WeatherProviderWithLaydown{}
}

func (w *WeatherProviderWithLaydown) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	resp, err := w.curr.GetWeather(city, ctx)

	if err != nil {
		if w.next == nil {
			return resp, err
		}

		return w.next.GetWeather(city, ctx)
	}

	return resp, nil
}

func (w *WeatherProviderWithLaydown) SetNext(next WeatherProvider) {}
