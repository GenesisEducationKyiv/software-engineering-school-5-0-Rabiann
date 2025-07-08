package providers_test

import (
	"testing"
	"time"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/Rabiann/weather-mailer/internal/notification/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest() *providers.RedisCache {
	cache := providers.NewRedisCache(nil)
	return cache
}

func TestCacheSetGet(t *testing.T) {
	cache := setupTest()
	defer func() {
		_ = cache.Clear()
	}()

	require.NoError(t, cache.Ping())

	key := "kyiv"

	weather := models.Weather{
		Temperature: 1,
		Humidity:    2,
		Description: "abc",
	}

	err := cache.Write(key, weather)
	assert.NoError(t, err)

	value, ok := cache.Read(key)
	assert.True(t, ok)
	require.Equal(t, weather, value)
}

func TestCacheSetGetExpired(t *testing.T) {
	cache := setupTest()
	defer func() {
		_ = cache.Clear()
	}()

	require.NoError(t, cache.Ping())
	key := "kyiv"

	weather := models.Weather{
		Temperature: 1,
		Humidity:    2,
		Description: "abc",
	}

	err := cache.SetExpired(key, weather, time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 10)
	value, ok := cache.Read(key)
	assert.False(t, ok)
	require.NotEqual(t, weather, value)
}
