package services

import (
	localConf "github.com/Rabiann/weather-mailer/pkg/go/api/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/producers"
)

type EmailService struct {
	producer *producers.EmailProducer
}

func NewEmailService(config *localConf.Configuration) (*EmailService, error) {
	producer, err := producers.NewEmailProducer(config)
	if err != nil {
		return nil, err
	}

	return &EmailService{producer}, nil
}

func (s *EmailService) SendConfirmationLetter(recipient string, confirmationUrl string) error {
	producer := s.producer

	msg := dto.SubscriptionConfirmationMsg{
		Email: recipient,
		Url:   confirmationUrl,
	}

	return producer.Publish_SendConfirmationLetter(msg)
}

func (s *EmailService) SendWeatherReport(subscriber dto.Subscriber, weather dto.Weather, unsubscriptionUrl string) error {
	producer := s.producer

	msg := dto.ReportMsg{
		UnsubscriptionUrl: unsubscriptionUrl,
		Subscriber:        subscriber,
		Weather:           weather,
	}
	return producer.Publish_SendWeatherReport(msg)
}
