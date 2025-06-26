package dto

import (
	"encoding/json"
	"fmt"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	WeatherApiRequestProvider struct {
		config *config.Configuration
		name   string
	}

	WeatherMapRequestProvider struct {
		config *config.Configuration
		name   string
	}

	WeatherStackRequestProvider struct {
		config *config.Configuration
		name   string
	}

	ApiRequestProvider interface {
		BuildUrl(city string) string
		BuildResponse([]byte) (models.Weather, error)
		Name() string
	}
)

func NewWeatherApiRequestProvider(config *config.Configuration, name string) *WeatherApiRequestProvider {
	return &WeatherApiRequestProvider{config, name}
}

func (p *WeatherApiRequestProvider) BuildUrl(city string) string {
	return fmt.Sprintf(p.config.WeatherApiAddress, p.config.WeatherApiKey, city)
}

func (p *WeatherApiRequestProvider) BuildResponse(body []byte) (models.Weather, error) {
	var weatherResponse models.WeatherApiResponse
	var weather models.Weather

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return weather, err
	}

	weather.Description = weatherResponse.Current.Condition.Text
	weather.Humidity = weatherResponse.Current.Humidity
	weather.Temperature = weatherResponse.Current.Temperature

	return weather, nil
}

func (p *WeatherApiRequestProvider) Name() string {
	return p.name
}

func NewWeatherStackRequestProvider(config *config.Configuration, name string) *WeatherApiRequestProvider {
	return &WeatherApiRequestProvider{config, name}
}

func (p *WeatherStackRequestProvider) BuildUrl(city string) string {
	return fmt.Sprintf(p.config.WeatherStackAddress, p.config.WeatherStackApi, city)
}

func (p *WeatherStackRequestProvider) BuildResponse(body []byte) (models.Weather, error) {
	var weatherResponse models.WeatherStackResponse
	var weather models.Weather

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return weather, err
	}

	weather.Description = weatherResponse.Current.Description
	weather.Humidity = weatherResponse.Current.Humidity
	weather.Temperature = weatherResponse.Current.Temperature

	return weather, nil
}

func (p *WeatherMapRequestProvider) Name() string {
	return p.name
}

func NewWeatherMapRequestProvider(config *config.Configuration, name string) *WeatherApiRequestProvider {
	return &WeatherApiRequestProvider{config, name}
}

func (p *WeatherMapRequestProvider) BuildUrl(city string) string {
	return fmt.Sprintf(p.config.WeatherMapAddress, city, p.config.WeatherMapApi)
}

func (p *WeatherMapRequestProvider) BuildResponse(body []byte) (models.Weather, error) {
	var weatherResponse models.WeatherMapResponse
	var weather models.Weather

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return weather, err
	}

	weather.Description = weatherResponse.Weather.Description
	weather.Humidity = weatherResponse.Main.Humidity
	weather.Temperature = weatherResponse.Main.Temperature

	return weather, nil
}

func (p *WeatherStackRequestProvider) Name() string {
	return p.name
}
