package providers

import (
	"sync"

	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	AsyncCache struct {
		cache map[string]models.Weather
		mu    sync.RWMutex
	}
)

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
