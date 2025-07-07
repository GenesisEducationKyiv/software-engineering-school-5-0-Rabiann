package notification

import "github.com/Rabiann/weather-mailer/internal/models"

type (
	CacheService struct {
		provider CacheProvider
	}

	CacheProvider interface {
		Read(key string) (models.Weather, bool)
		Write(key string, value models.Weather) error
	}
)

func NewCacheService(provider CacheProvider) *CacheService {
	return &CacheService{
		provider: provider,
	}
}

func (c *CacheService) Read(key string) (models.Weather, bool) {
	return c.provider.Read(key)
}

func (c *CacheService) Write(key string, value models.Weather) error {
	return c.provider.Write(key, value)
}
