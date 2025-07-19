package weather

import (
	"context"

	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	WeatherProviderWithLaydown struct {
		curr WeatherProviderer
		next LaydownWeatherProvider
	}

	LaydownWeatherProvider interface {
		GetWeather(city string, ctx context.Context) (models.Weather, error)
		Add(next WeatherProviderer)
	}

	WeatherProviderer interface {
		GetWeather(city string, ctx context.Context) (models.Weather, error)
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

func (w *WeatherProviderWithLaydown) Add(next WeatherProviderer) {
	if w.curr == nil {
		w.curr = next
		return
	}

	if w.next == nil {
		w.next = &WeatherProviderWithLaydown{}
		w.next.Add(next)
		return
	}

	w.next.Add(next)
}
