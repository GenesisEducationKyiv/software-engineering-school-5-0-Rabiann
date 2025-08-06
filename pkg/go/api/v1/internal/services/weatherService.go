package services

import (
	"context"

	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
	pWeather "github.com/Rabiann/weather-mailer/protos/weather"
	"google.golang.org/grpc"
)

type WeatherService struct {
	weatherClient pWeather.WeatherServiceClient
}

func NewWeatherService(conn *grpc.ClientConn) *WeatherService {
	c := pWeather.NewWeatherServiceClient(conn)
	return &WeatherService{weatherClient: c}
}

func (w *WeatherService) GetWeather(city string, ctx context.Context) (dto.Weather, error) {
	var weather dto.Weather
	weatherClient := w.weatherClient
	reply, err := weatherClient.GetWeather(ctx, &pWeather.WeatherRequest{City: city})
	if err != nil {
		return weather, err
	}

	return dto.Weather{
		Temperature: float64(reply.Temperature),
		Humidity:    float64(reply.Humidity),
		Description: reply.Description,
	}, nil
}
