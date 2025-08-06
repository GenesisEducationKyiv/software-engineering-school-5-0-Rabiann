package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
	pUser "github.com/Rabiann/weather-mailer/protos/user"
	"github.com/go-co-op/gocron/v2"
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
		userClient          pUser.UserServiceClient
	}

	MailingService interface {
		SendWeatherReport(subscriber dto.Subscriber, weather dto.Weather, unsubscriptionUrl string) error
	}

	SubscriptionService interface {
		GetActiveSubscriptions(string, context.Context) ([]dto.Subscription, error)
	}

	WeatherService interface {
		GetWeather(string, context.Context) (dto.Weather, error)
	}
)

func NewNotifier(weatherService WeatherService, subscriptionService SubscriptionService, mailingService MailingService) Notifier {
	return Notifier{
		weatherService:      weatherService,
		subscriptionService: subscriptionService,
		mailingService:      mailingService,
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
		go func(dto.Subscription) {
			weather, err := n.weatherService.GetWeather(sub.City, ctx_)

			tokenRequest := pUser.CreateTokenRequest{
				SubscriptionId: int32(sub.Id),
			}
			token, err := n.userClient.CreateToken(ctx_, &tokenRequest)
			if err != nil {
				return
			}

			url := fmt.Sprintf("%s/api/unsubscribe/%s", baseUrl, token)

			sub := dto.Subscriber{
				Recipient: sub.Email,
				Period:    per,
				City:      sub.City,
			}
			_ = n.mailingService.SendWeatherReport(sub, weather, url)
		}(sub)
	}

	return nil
}
