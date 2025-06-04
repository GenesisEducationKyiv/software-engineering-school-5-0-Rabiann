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

	config.baseUrl = os.Getenv("BASE_URL")
	config.sendgridApiKey = os.Getenv("SENDGRID_API_KEY")
	config.weatherApiKey = os.Getenv("WEATHER_API_KEY")
	config.senderMail = os.Getenv("SENDER_MAIL")
	config.weatherApiAddress = os.Getenv("WEATHER_API_ADDR")
	config.port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	return &config
}
