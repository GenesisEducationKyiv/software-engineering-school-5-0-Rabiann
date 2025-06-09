package cmd

import (
	"net/http"
	"strconv"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/controllers"
	"github.com/Rabiann/weather-mailer/internal/notification"
	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/Rabiann/weather-mailer/internal/services/models"
	"github.com/gin-gonic/gin"
)

type App struct{}

func (a *App) Run() error {
	configuration, err := config.LoadEnvironment()
	if err != nil {
		return err
	}

	db := models.ConnectToDatabase()

	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Token{}); err != nil {
		return err
	}

	weatherService := services.WeatherService{Config: configuration}
	subscriptionService := services.SubscriptionService{Db: db}
	tokenService := services.TokenService{Db: db}
	emailService, err := services.NewMailingService(configuration)
	if err != nil {
		return err
	}

	notifier := notification.NewNotifier(weatherService, subscriptionService, emailService, tokenService, configuration)
	go notifier.RunNotifier()

	weatherController := controllers.WeatherController{WeatherService: weatherService}
	subscriptionController := controllers.SubscriptionController{SubscriptionService: subscriptionService, TokenService: tokenService, EmailService: emailService, BaseUrl: configuration.BaseUrl}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/favicon.ico", "./static/weather.ico")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "subscriptions.html", gin.H{})
	})

	api := router.Group("/api")
	{
		api.GET("/weather", weatherController.GetWeather)
		api.POST("/subscribe", subscriptionController.Subscribe)
		api.GET("/confirm/:token", subscriptionController.Confirm)
		api.GET("/unsubscribe/:token", subscriptionController.Unsubscribe)
	}

	if err := router.Run(strconv.Itoa(configuration.Port)); err != nil {
		return err
	}

	return nil
}
