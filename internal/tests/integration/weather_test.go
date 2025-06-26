package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/controllers"
	"github.com/Rabiann/weather-mailer/internal/dto"
	"github.com/Rabiann/weather-mailer/internal/external/weather"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupWeatherServer(response any) *httptest.Server {
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

	apiProvider := dto.NewWeatherApiRequestProvider(configuration, "weatherapi.org")
	weatherProvider := weather.NewWeatherProvider(configuration, apiProvider)
	weatherService := services.NewWeatherService(weatherProvider)
	weatherController := controllers.NewWeatherController(weatherService)
	return weatherController
}

func TestGetWeatherApi(t *testing.T) {
	weather := models.WeatherApiResponse{
		Current: struct {
			Temperature float64 "json:\"temp_c\""
			Humidity    float64 "json:\"humidity\""
			Condition   struct {
				Text string "json:\"text\""
			} "json:\"condition\""
		}{},
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
	assert.Equal(t, weatherResponse.Temperature, weather.Current.Temperature)
	assert.Equal(t, weatherResponse.Humidity, weather.Current.Humidity)
	assert.Equal(t, weatherResponse.Description, weather.Current.Condition.Text)
}

func TestGetWeatherMap(t *testing.T) {
	weather := models.WeatherMapResponse{
		Weather: struct {
			Description string "json:\"description\""
		}{},
		Main: struct {
			Temperature float64 "json:\"temp\""
			Humidity    float64 "json:\"humidity\""
		}{},
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
	assert.Equal(t, weatherResponse.Temperature, weather.Main.Temperature)
	assert.Equal(t, weatherResponse.Humidity, weather.Main.Humidity)
	assert.Equal(t, weatherResponse.Description, weather.Weather.Description)
}

func TestGetWeatherStack(t *testing.T) {
	weather := models.WeatherStackResponse{
		Current: struct {
			Temperature float64 "json:\"temperature\""
			Description string  "json:\"weather_description\""
			Humidity    float64 "json:\"humidity\""
		}{},
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
	assert.Equal(t, weatherResponse.Temperature, weather.Current.Temperature)
	assert.Equal(t, weatherResponse.Humidity, weather.Current.Humidity)
	assert.Equal(t, weatherResponse.Description, weather.Current.Description)
}
