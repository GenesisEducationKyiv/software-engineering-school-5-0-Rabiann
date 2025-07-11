package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/models"
)

type WeatherProvider struct {
	config *config.Configuration
	client *http.Client
}

func NewWeatherProvider(config *config.Configuration) *WeatherProvider {
	return &WeatherProvider{config, &http.Client{}}
}

func (w *WeatherProvider) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	var weather models.Weather
	var weatherResponse models.WeatherResponse
	url := fmt.Sprintf(w.config.WeatherApiAddress, w.config.WeatherApiKey, city)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return weather, err
	}

	weather.Description = weatherResponse.Text
	weather.Humidity = weatherResponse.Humidity
	weather.Temperature = weatherResponse.Temperature

	return weather, nil
}
