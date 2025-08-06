package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"unicode"

	"github.com/joho/godotenv"
)

// type Configuration struct {
// 	BaseUrl             string
// 	SendgridApiKey      string
// 	WeatherApiKey       string
// 	SenderMail          string
// 	WeatherApiAddress   string `env:"WEATHER_API_ADDR"`
// 	Port                string
// 	MailTimeout         int
// 	WeatherMapAddress   string `env:"WEATHER_MAP_ADDR"`
// 	WeatherMapApi       string
// 	WeatherStackAddress string `env:"WEATHER_STACK_ADDR"`
// 	WeatherStackApi     string `env:"WEATHER_STACK_KEY"`
// 	RedisAddr           string
// 	LocalCacheFallback  string
// }

func FromCamelCaseToUpperCase(text string) string {
	result := ""

	for index, chr := range text {
		if index > 0 && unicode.IsUpper(chr) {
			result += "_"
		}

		result += string(unicode.ToUpper(chr))
	}

	return result
}

func LoadEnvironment(configuration any) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("`.env` is not found. Using user environment")
	}

	v := reflect.ValueOf(configuration)
	v = v.Elem()
	type_ := v.Type()

	fieldCount := type_.NumField()

	for i := range fieldCount {
		field := type_.Field(i)
		name := field.Name

		envName := field.Tag.Get("env")
		if envName == "" {
			envName = FromCamelCaseToUpperCase(name)
		}

		env := os.Getenv(envName)
		if env == "" {
			return fmt.Errorf("`%s` not set", envName)
		}

		switch field.Type.Kind() {
		case reflect.Int:
			value, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("`%s` should be Int", envName)
			}

			v.Field(i).SetInt(int64(value))
		default:
			v.Field(i).SetString(env)
		}
	}

	return nil
}
