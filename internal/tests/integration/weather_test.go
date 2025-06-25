package integration_test

import (
	"encoding/json"
	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/controllers"
	"github.com/Rabiann/weather-mailer/internal/external"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupWeatherServer(response models.WeatherResponse) *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		weatherResponse := response
		body, err := json.Marshal(weatherResponse)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))

	return mockServer
}

func setupWeatherTest(mockServerUrl string) controllers.WeatherController {
	configuration := &config.Configuration{
		WeatherApiKey:     "testApikey",
		WeatherApiAddress: mockServerUrl + "/weather?key=%s&q=%s&aqi=no",
	}

	weatherProvider := external.NewWeatherProvider(configuration)
	weatherService := services.NewWeatherService(weatherProvider)
	weatherController := controllers.NewWeatherController(weatherService)
	return weatherController
}

func TestGetWeather(t *testing.T) {
	weather := models.WeatherResponse{
		Current: models.Current{
			Temperature: 1,
			Humidity:    1,
			Condition: models.Condition{
				Text: "cold and dry",
			},
		},
	}

	server := setupWeatherServer(weather)
	defer server.Close()
	weatherController := setupWeatherTest(server.URL)

	router := gin.Default()

	api := router.Group("/api")
	api.GET("/weather", weatherController.GetWeather)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/weather?key=testApikey&city=kyiv&aqi=no", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var weatherResponse models.Weather
	err := json.Unmarshal(w.Body.Bytes(), &weatherResponse)
	assert.NoError(t, err)
	assert.Equal(t, weatherResponse.Temperature, weather.Temperature)
	assert.Equal(t, weatherResponse.Humidity, weather.Humidity)
	assert.Equal(t, weatherResponse.Description, weather.Condition.Text)
}
