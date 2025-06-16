package notification

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/Rabiann/weather-mailer/internal/services/models"
	"github.com/go-co-op/gocron/v2"
)

const Day = time.Hour * 24

type Period int

const (
	Hourly Period = iota
	Daily
)

type AsyncCache struct {
	cache map[string]services.Weather
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
	cache := make(map[string]services.Weather)
	return AsyncCache{
		cache: cache,
		mu:    sync.RWMutex{},
	}
}

func (c *AsyncCache) Read(key string) (services.Weather, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wth, ok := c.cache[key]
	return wth, ok
}

func (c *AsyncCache) Write(key string, value services.Weather) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

type Notifier struct {
	weatherService      *services.WeatherService
	subscriptionService *services.SubscriptionService
	mailingService      *services.MailingService
	tokenService        *services.TokenService
}

func NewNotifier(weatherService *services.WeatherService, subscriptionService *services.SubscriptionService, mailingService *services.MailingService, tokenService *services.TokenService) Notifier {
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

func (n Notifier) RunSendingPipeline(period Period, baseUrl string) {
	var subscribers []models.Subscription
	var per string
	var err error

	cache := NewAsyncCache()
	semaphore := NewSemaphore(10)

	if period == Daily {
		per = "daily"
	} else {
		per = "hourly"
	}

	result := n.subscriptionService.Db.Where("frequency = ? and confirmed = true", per).Find(&subscribers)
	if result.Error != nil {
		panic(result.Error)
	}

	for _, sub := range subscribers {
		semaphore.Acquire()
		go func(models.Subscription) {
			defer semaphore.Release()
			city := strings.ToLower(sub.City)
			weather, ok := cache.Read(city)

			if !ok {
				weather, err = n.weatherService.GetWeather(city)
				if err != nil {
					return
				}

				cache.Write(sub.City, weather)
			}

			token, err := n.tokenService.CreateToken(sub.ID)
			if err != nil {
				return
			}

			url := fmt.Sprintf("%s/api/unsubscribe/%s", baseUrl, token)

			sub := services.Subscriber{
				Recipient: sub.Email,
				Period:    per,
				City:      sub.City,
			}
			_ = n.mailingService.SendWeatherReport(&sub, &weather, url)
		}(sub)
	}
}
