package notification

import (
	"context"
	"fmt"
	"strings"
	"sync"
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

type AsyncCache struct {
	cache map[string]models.Weather
	mu    sync.RWMutex
}

type Semaphore struct {
	c chan struct{}
}

func (s *Semaphore) Acquire() {
	s.c <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.c
}

func NewSemaphore(wCount int) Semaphore {
	c := make(chan struct{}, wCount)
	return Semaphore{c: c}
}

func NewAsyncCache() AsyncCache {
	cache := make(map[string]models.Weather)
	return AsyncCache{
		cache: cache,
		mu:    sync.RWMutex{},
	}
}

func (c *AsyncCache) Read(key string) (models.Weather, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wth, ok := c.cache[key]
	return wth, ok
}

func (c *AsyncCache) Write(key string, value models.Weather) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

type (
	Notifier struct {
		weatherService      WeatherService
		subscriptionService SubscriptionService
		mailingService      MailingService
		tokenService        TokenService
	}

	MailingService interface {
		SendWeatherReport(*models.Subscriber, *models.Weather, string) error
	}

	TokenService interface {
		CreateToken(uint, context.Context, context.CancelFunc) (uuid.UUID, error)
	}

	SubscriptionService interface {
		GetActiveSubscriptions(string, context.Context, context.CancelFunc) ([]models.Subscription, error)
	}

	WeatherService interface {
		GetWeather(string, context.Context, context.CancelFunc) (models.Weather, error)
	}
)

func NewNotifier(weatherService WeatherService, subscriptionService SubscriptionService, mailingService MailingService, tokenService TokenService) Notifier {
	return Notifier{
		weatherService:      weatherService,
		subscriptionService: subscriptionService,
		mailingService:      mailingService,
		tokenService:        tokenService,
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

func (n Notifier) RunSendingPipeline(period Period, baseUrl string) error {
	var per string
	var err error

	cache := NewAsyncCache()
	semaphore := NewSemaphore(10)

	if period == Daily {
		per = "daily"
	} else {
		per = "hourly"
	}

	ctx_, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	subscribers, err := n.subscriptionService.GetActiveSubscriptions(per, ctx_, cancel)
	if err != nil {
		return err
	}

	for _, sub := range subscribers {
		semaphore.Acquire()
		go func(models.Subscription) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer semaphore.Release()
			city := strings.ToLower(sub.City)
			weather, ok := cache.Read(city)

			if !ok {
				weather, err = n.weatherService.GetWeather(city, ctx, cancel)
				if err != nil {
					return
				}

				cache.Write(sub.City, weather)
			}

			token, err := n.tokenService.CreateToken(sub.ID, ctx_, cancel)
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
