package weather_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Rabiann/weather-mailer/internal/external/weather"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	MockProvider1 struct {
		mock.Mock
	}

	MockProvider2 struct {
		mock.Mock
	}

	MockProvider3 struct {
		mock.Mock
	}
)

func (p *MockProvider1) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	res := p.Called(city, ctx)
	return res.Get(0).(models.Weather), res.Error(1)

}

func (p *MockProvider2) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	res := p.Called(city, ctx)
	return res.Get(0).(models.Weather), res.Error(1)

}

func (p *MockProvider3) GetWeather(city string, ctx context.Context) (models.Weather, error) {
	res := p.Called(city, ctx)
	return res.Get(0).(models.Weather), res.Error(1)

}

func TestWeatherProviderWithLaydownAllSuccess(t *testing.T) {
	provider1 := new(MockProvider1)
	provider2 := new(MockProvider2)
	provider3 := new(MockProvider3)
	weather1 := models.Weather{Temperature: 1, Humidity: 2, Description: "TEXT1"}
	weather2 := models.Weather{Temperature: 2, Humidity: 2, Description: "TEXT2"}
	weather3 := models.Weather{Temperature: 3, Humidity: 2, Description: "TEXT3"}
	provider1.On("GetWeather", "city", context.TODO()).Return(weather1, nil)
	provider2.On("GetWeather", "city", context.TODO()).Return(weather2, nil)
	provider3.On("GetWeather", "city", context.TODO()).Return(weather3, nil)

	laydown := weather.WeatherProviderWithLaydown{}
	laydown.Add(provider1)
	laydown.Add(provider2)
	laydown.Add(provider3)

	res, err := laydown.GetWeather("city", context.TODO())
	assert.NoError(t, err)

	require.Equal(t, weather1.Temperature, res.Temperature)
}

func TestWeatherProviderWithLaydownFirstFailed(t *testing.T) {
	provider1 := new(MockProvider1)
	provider2 := new(MockProvider2)
	provider3 := new(MockProvider3)
	weather1 := models.Weather{Temperature: 1, Humidity: 2, Description: "TEXT1"}
	weather2 := models.Weather{Temperature: 2, Humidity: 2, Description: "TEXT2"}
	weather3 := models.Weather{Temperature: 3, Humidity: 2, Description: "TEXT3"}
	provider1.On("GetWeather", "city", context.TODO()).Return(weather1, errors.New("ERROR1"))
	provider2.On("GetWeather", "city", context.TODO()).Return(weather2, nil)
	provider3.On("GetWeather", "city", context.TODO()).Return(weather3, nil)

	laydown := weather.WeatherProviderWithLaydown{}
	laydown.Add(provider1)
	laydown.Add(provider2)
	laydown.Add(provider3)

	res, err := laydown.GetWeather("city", context.TODO())
	assert.NoError(t, err)

	require.Equal(t, weather2.Temperature, res.Temperature)
}

func TestWeatherProviderWithLaydownSecondFailed(t *testing.T) {
	provider1 := new(MockProvider1)
	provider2 := new(MockProvider2)
	provider3 := new(MockProvider3)
	weather1 := models.Weather{Temperature: 1, Humidity: 2, Description: "TEXT1"}
	weather2 := models.Weather{Temperature: 2, Humidity: 2, Description: "TEXT2"}
	weather3 := models.Weather{Temperature: 3, Humidity: 2, Description: "TEXT3"}
	provider1.On("GetWeather", "city", context.TODO()).Return(weather1, errors.New("ERROR1"))
	provider2.On("GetWeather", "city", context.TODO()).Return(weather2, errors.New("ERROR2"))
	provider3.On("GetWeather", "city", context.TODO()).Return(weather3, nil)

	laydown := weather.WeatherProviderWithLaydown{}
	laydown.Add(provider1)
	laydown.Add(provider2)
	laydown.Add(provider3)

	res, err := laydown.GetWeather("city", context.TODO())
	assert.NoError(t, err)

	require.Equal(t, weather3.Temperature, res.Temperature)
}
