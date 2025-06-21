package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/controllers"
	"github.com/Rabiann/weather-mailer/internal/external"
	"github.com/Rabiann/weather-mailer/internal/notification"
	"github.com/Rabiann/weather-mailer/internal/persistance"
	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct{}

func bootstrapDatabase() (*gorm.DB, error) {
	db := persistance.ConnectToDatabase()

	if err := persistance.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func (a *App) Run() error {
	configuration, err := config.LoadEnvironment()
	if err != nil {
		return err
	}

	db, err := bootstrapDatabase()
	if err != nil {
		return err
	}

	subscriptionRepository := persistance.NewSubscriptionRepository(db)
	tokenRepository := persistance.NewTokenRepository(db)
	weatherProvider := external.NewWeatherProvider(configuration)
	mailingProvider := external.NewMailingProvider(configuration)

	weatherService := services.NewWeatherService(weatherProvider)
	subscriptionDataService := services.NewSubscriptionService(subscriptionRepository)
	tokenService := services.NewTokenService(tokenRepository)
	emailService, err := services.NewMailingService(mailingProvider, configuration)
	if err != nil {
		return err
	}

	subscriptionService := services.NewSubscriptionBusinessService(subscriptionDataService, tokenService, emailService, configuration.BaseUrl)
	notifier := notification.NewNotifier(weatherService, subscriptionDataService, emailService, tokenService)
	go notifier.RunNotifier(configuration.BaseUrl)

	weatherController := controllers.NewWeatherController(weatherService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)
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

	srv := &http.Server{
		Addr:    ":" + configuration.Port,
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown server.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Print("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout 5 seconds")
	log.Printf("server exiting")

	return nil
}
