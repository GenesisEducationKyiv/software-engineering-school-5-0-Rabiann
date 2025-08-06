package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rabiann/weather-mailer/lib/go/config"
	"github.com/Rabiann/weather-mailer/lib/go/logger"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/api/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/controllers"
	notification "github.com/Rabiann/weather-mailer/pkg/go/api/internal/notifier"
	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

type App struct{}

func (a *App) Run() error {
	configuration := localConf.Configuration{}
	err := config.LoadEnvironment(&configuration)
	if err != nil {
		return err
	}

	os.Mkdir("logs", 0700)
	logger.SetupLogger("logs/app.log")

	userConn, err := grpc.NewClient("localhost:5555", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Can't connect to user microservice")
	}

	weatherConn, err := grpc.NewClient("localhost:5556", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Can't connect to weather microservice")
	}

	subscriptionService := services.NewSubscriptionBusinessService(userConn, nil, "123")
	weatherService := services.NewWeatherService(weatherConn)
	weatherController := controllers.NewWeatherController(weatherService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)
	mailingService, err := services.NewEmailService(&configuration)
	if err != nil {
		return err
	}

	notifier := notification.NewNotifier(weatherService, subscriptionService, mailingService)

	go notifier.RunNotifier(configuration.BaseUrl)

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
