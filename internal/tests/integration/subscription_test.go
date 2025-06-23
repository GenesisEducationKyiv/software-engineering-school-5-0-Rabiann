package integration_test

import (
	"fmt"
	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/controllers"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/Rabiann/weather-mailer/internal/persistance"
	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func setupSubscriptionTest() (*gin.Engine, *gorm.DB) {
	configuration := &config.Configuration{
		WeatherApiKey: "testApikey",
		BaseUrl:       "baseurl",
	}

	db, err := persistance.SetupInMemoryDb()
	if err != nil {
		panic(err)
	}

	subscriptionRepository := persistance.NewSubscriptionRepository(db)
	tokenRepository := persistance.NewTokenRepository(db)

	subscriptionDataService := services.NewSubscriptionService(subscriptionRepository)
	tokenService := services.NewTokenService(tokenRepository)
	emailService := new(MockMailService)
	emailService.On("SendConfirmationLetter", mock.Anything, mock.Anything).Return(nil)
	emailService.On("SendWeatherReport", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	subscriptionService := services.NewSubscriptionBusinessService(subscriptionDataService, tokenService, emailService, configuration.BaseUrl)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	router := gin.Default()
	router.LoadHTMLGlob("../../../templates/*")

	api := router.Group("/api")
	api.POST("/subscribe", subscriptionController.Subscribe)
	api.GET("/confirm/:token", subscriptionController.Confirm)
	api.GET("/unsubscribe/:token", subscriptionController.Unsubscribe)

	return router, db
}

func TestSubscribe(t *testing.T) {
	router, _ := setupSubscriptionTest()
	w := httptest.NewRecorder()

	subscription := models.Subscription{
		Email:     "vasya@mail.com",
		City:      "kyiv",
		Frequency: "daily",
	}

	form := url.Values{}
	form.Add("email", subscription.Email)
	form.Add("city", subscription.City)
	form.Add("period", subscription.Frequency)

	req, err := http.NewRequest("POST", "/api/subscribe", strings.NewReader(form.Encode()))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestConfirm(t *testing.T) {
	router, db := setupSubscriptionTest()

	token := uuid.New()

	db.Create(&models.Subscription{
		ID:        1,
		Email:     "vasya@mail.com",
		City:      "kyiv",
		Frequency: "daily",
		Confirmed: false,
	})

	db.Create(&models.Token{
		ID:             token,
		Expires:        time.Now().AddDate(1000, 0, 0),
		SubscriptionID: 1,
	})

	w := httptest.NewRecorder()
	reqUrl := fmt.Sprintf("/api/confirm/%s", token)

	req, err := http.NewRequest("GET", reqUrl, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUnsubscribe(t *testing.T) {
	router, db := setupSubscriptionTest()

	token := uuid.New()

	db.Create(&models.Subscription{
		ID:        1,
		Email:     "vasya@mail.com",
		City:      "kyiv",
		Frequency: "daily",
		Confirmed: true,
	})

	db.Create(&models.Token{
		ID:             token,
		Expires:        time.Now().AddDate(1000, 0, 0),
		SubscriptionID: 1,
	})

	w := httptest.NewRecorder()
	reqUrl := fmt.Sprintf("/api/unsubscribe/%s", token)

	req, err := http.NewRequest("GET", reqUrl, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
