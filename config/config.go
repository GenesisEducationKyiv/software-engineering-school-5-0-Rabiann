package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	BaseUrl           string
	SendgridApiKey    string
	WeatherApiKey     string
	SenderMail        string
	WeatherApiAddress string
	Port              int
}

func LoadEnvironment() *Configuration {
	var config Configuration
	var err error
	_ = godotenv.Load(".env")

	config.BaseUrl = os.Getenv("BASE_URL")
	config.SendgridApiKey = os.Getenv("SENDGRID_API_KEY")
	config.WeatherApiKey = os.Getenv("WEATHER_API_KEY")
	config.SenderMail = os.Getenv("SENDER_MAIL")
	config.WeatherApiAddress = os.Getenv("WEATHER_API_ADDR")
	config.Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	return &config
}
