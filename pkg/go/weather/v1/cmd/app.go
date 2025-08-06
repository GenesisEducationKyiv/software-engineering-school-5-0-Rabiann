package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Rabiann/weather-mailer/lib/go/config"
	"github.com/Rabiann/weather-mailer/lib/go/logger"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/weather/internal/config"
	weather "github.com/Rabiann/weather-mailer/pkg/go/weather/internal/external/weather"
	"github.com/Rabiann/weather-mailer/pkg/go/weather/internal/providers"
	"github.com/Rabiann/weather-mailer/pkg/go/weather/internal/services"
	pb "github.com/Rabiann/weather-mailer/protos/weather"
	"google.golang.org/grpc"
)

type App struct{}

type Server struct {
	pb.UnimplementedWeatherServiceServer
	weatherService *services.WeatherService
}

func (s *Server) GetWeather(ctx context.Context, in *pb.WeatherRequest) (*pb.WeatherReply, error) {
	log.Printf("Get weather")
	weatherService := s.weatherService
	city := in.City

	weather, err := weatherService.GetWeather(city, ctx)
	if err != nil {
		return nil, err
	}

	reply := &pb.WeatherReply{
		Temperature: float32(weather.Temperature),
		Humidity:    float32(weather.Humidity),
		Description: weather.Description,
	}

	return reply, nil
}

func (a *App) Run() error {
	configuration := &localConf.Configuration{}
	err := config.LoadEnvironment(configuration)
	if err != nil {
		return err
	}

	os.Mkdir("logs", 0700)
	logger.SetupLogger("logs/app.log")
	cacheProvider := providers.NewRedisCache(configuration)
	cacheService := services.NewCacheService(cacheProvider)
	weatherApiProvider := weather.NewWeatherProviderLogger(weather.NewWeatherProvider(configuration, weather.NewWeatherApiRequestProvider(configuration, "weatherapi.org")))
	weatherMapProvider := weather.NewWeatherProviderLogger(weather.NewWeatherProvider(configuration, weather.NewWeatherMapRequestProvider(configuration, "openweathermap.org")))
	weatherStackProvider := weather.NewWeatherProviderLogger(weather.NewWeatherProvider(configuration, weather.NewWeatherStackRequestProvider(configuration, "weatherstack.org")))
	weatherProvider := weather.NewWeatherProviderWithLaydown()
	weatherProvider.Add(weatherApiProvider)
	weatherProvider.Add(weatherMapProvider)
	weatherProvider.Add(weatherStackProvider)
	weatherService := services.NewWeatherService(weatherProvider, cacheService)

	server := &Server{
		weatherService: weatherService,
	}

	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5556))
	pb.RegisterWeatherServiceServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

	return nil
}
