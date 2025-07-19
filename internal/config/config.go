package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	BaseUrl             string
	SendgridApiKey      string
	WeatherApiKey       string
	SenderMail          string
	WeatherApiAddress   string
	Port                string
	MailTimeout         int
	WeatherMapAddress   string
	WeatherMapApi       string
	WeatherStackAddress string
	WeatherStackApi     string
}

func LoadEnvironment() (*Configuration, error) {
	var config Configuration
	var err error
	if err = godotenv.Load(".env"); err != nil {
		log.Fatal("`.env` is not found. Using user environment")
	}

	config.BaseUrl = os.Getenv("BASE_URL")
	if config.BaseUrl == "" {
		return nil, errors.New("`BASE_URL` is not set")
	}

	config.SendgridApiKey = os.Getenv("SENDGRID_API_KEY")
	if config.BaseUrl == "" {
		return nil, errors.New("`SENDGRID_API_KEY` is not set")
	}

	config.WeatherApiKey = os.Getenv("WEATHER_API_KEY")
	if config.BaseUrl == "" {
		return nil, errors.New("`WEATHER_API_KEY` is not set")
	}

	config.WeatherMapApi = os.Getenv("WEATHER_MAP_API")
	if config.BaseUrl == "" {
		return nil, errors.New("`WEATHER_MAP_API` is not set")
	}

	config.WeatherStackApi = os.Getenv("WEATHER_STACK_API")
	if config.BaseUrl == "" {
		return nil, errors.New("WEATHER_STACK_API is not set")
	}

	config.SenderMail = os.Getenv("SENDER_MAIL")
	if config.SenderMail == "" {
		return nil, errors.New("`SENDER_MAIL` is not set")
	}

	config.WeatherApiAddress = os.Getenv("WEATHER_API_ADDR")
	if config.WeatherApiAddress == "" {
		return nil, errors.New("`WEATHER_API_ADDR` is not set")
	}

	config.WeatherMapAddress = os.Getenv("WEATHER_MAP_ADDR")
	if config.WeatherMapAddress == "" {
		return nil, errors.New("`WEATHER_MAP_ADDR` is not set")
	}

	config.WeatherStackAddress = os.Getenv("WEATHER_STACK_ADDR")
	if config.WeatherStackAddress == "" {
		return nil, errors.New("`WEATHER_STACK_ADDR` is not set")
	}

	mailTimeout := os.Getenv("MAIL_TIMEOUT")
	if mailTimeout == "" {
		return nil, errors.New("`MAIL_TIMEOUT` is not set")
	}

	config.MailTimeout, err = strconv.Atoi(mailTimeout)
	if err != nil {
		return nil, errors.New("`MAIL_TIMEOUT` should be valid integer")
	}

	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		return nil, errors.New("`PORT` is not set")
	}

	return &config, nil
}
