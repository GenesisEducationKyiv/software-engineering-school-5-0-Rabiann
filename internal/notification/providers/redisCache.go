package providers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/redis/go-redis/v9"
)

type (
	RedisCache struct {
		client *redis.Client
	}
)

func NewRedisCache(config *config.Configuration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) Read(key string) (models.Weather, bool) {
	ctx := context.Background()
	var weather models.Weather

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return weather, false
	}

	if err := json.Unmarshal([]byte(val), &weather); err != nil {
		return weather, false
	}

	return weather, true
}

func (c *RedisCache) Write(key string, value models.Weather) error {
	return c.SetExpired(key, value, time.Minute*10)
}

func (c *RedisCache) Clear() error {
	ctx := context.Background()
	err := c.client.FlushDB(ctx).Err()
	return err
}

func (c *RedisCache) Ping() error {
	ctx := context.Background()
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) SetExpired(key string, value models.Weather, expirationPeriod time.Duration) error {
	ctx := context.Background()
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = c.client.Set(ctx, key, string(body), expirationPeriod).Result()
	if err != nil {
		return err
	}

	return nil
}
