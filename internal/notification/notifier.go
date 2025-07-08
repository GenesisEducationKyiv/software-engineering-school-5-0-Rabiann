package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

const Day = time.Hour * 24

type Period int

const (
	Hourly Period = iota
	Daily
)

type (
	Notifier struct {
		weatherService      WeatherService
		subscriptionService SubscriptionService
		mailingService      MailingService
		tokenService        TokenService
		cacheService        CacheServicer
	}

	MailingService interface {
		SendWeatherReport(*models.Subscriber, *models.Weather, string) error
	}

	TokenService interface {
		CreateToken(uint, context.Context) (uuid.UUID, error)
	}

	SubscriptionService interface {
		GetActiveSubscriptions(string, context.Context) ([]models.Subscription, error)
	}

	WeatherService interface {
		GetWeather(string, context.Context) (models.Weather, error)
	}

	CacheServicer interface {
		Read(key string) (models.Weather, bool)
		Write(key string, value models.Weather) error
	}
)

func NewNotifier(weatherService WeatherService, subscriptionService SubscriptionService, mailingService MailingService, tokenService TokenService, cacheService CacheServicer) Notifier {
	return Notifier{
		weatherService:      weatherService,
		subscriptionService: subscriptionService,
		mailingService:      mailingService,
		tokenService:        tokenService,
		cacheService:        cacheService,
	}
}

func (n Notifier) RunNotifier(baseUrl string) {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			Day,
		),
		gocron.NewTask(
			n.RunSendingPipeline,
			Daily,
			baseUrl,
		),
	)

	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			time.Hour,
		),
		gocron.NewTask(
			n.RunSendingPipeline,
			Hourly,
			baseUrl,
		),
	)

	if err != nil {
		panic(err)
	}

	s.Start()

	// block thread, run scheduler infinitely
	select {}
}

func (n *Notifier) RunSendingPipeline(period Period, baseUrl string) error {
	cache := n.cacheService
	var per string
	var err error

	if period == Daily {
		per = "daily"
	} else {
		per = "hourly"
	}

	ctx_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	subscribers, err := n.subscriptionService.GetActiveSubscriptions(per, ctx_)
	if err != nil {
		return err
	}

	for _, sub := range subscribers {
		go func(models.Subscription) {
			city := strings.ToLower(sub.City)
			weather, ok := cache.Read(city)

			if !ok {
				weather, err = n.weatherService.GetWeather(city, ctx_)
				if err != nil {
					return
				}

				_ = cache.Write(sub.City, weather)
			}

			token, err := n.tokenService.CreateToken(sub.ID, ctx_)
			if err != nil {
				return
			}

			url := fmt.Sprintf("%s/api/unsubscribe/%s", baseUrl, token)

			sub := models.Subscriber{
				Recipient: sub.Email,
				Period:    per,
				City:      sub.City,
			}
			_ = n.mailingService.SendWeatherReport(&sub, &weather, url)
		}(sub)
	}

	return nil
}
