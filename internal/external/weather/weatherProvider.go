package weather

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	WeatherProvider struct {
		config      *config.Configuration
		client      *http.Client
		apiProvider ApiProvider
	}

	ApiProvider interface {
		BuildUrl(city string) string
		BuildResponse([]byte) (models.Weather, error)
		Name() string
	}
)

func NewWeatherProvider(config *config.Configuration, provider ApiProvider) *WeatherProvider {
	return &WeatherProvider{config, &http.Client{}, provider}
}

func (w *WeatherProvider) GetWeather(city string, ctx_ context.Context) (models.Weather, error) {
	var weather models.Weather
	url := w.apiProvider.BuildUrl(city)

	req, err := http.NewRequestWithContext(ctx_, "GET", url, nil)
	if err != nil {
		return weather, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return weather, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return weather, fmt.Errorf("city `%s` not exists", city)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather, err
	}

	weather, err = w.apiProvider.BuildResponse(body)
	return weather, err
}

func (w *WeatherProvider) Name() string {
	return w.apiProvider.Name()
}
