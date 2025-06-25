package integration_test

import (
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/stretchr/testify/mock"
)

type (
	MockMailService struct {
		mock.Mock
	}
)

func (s *MockMailService) SendConfirmationLetter(recipient string, confirmationUrl string) error {
	s.Called(recipient, confirmationUrl)
	return nil
}

func (s *MockMailService) SendWeatherReport(subscriber *models.Subscriber, weather *models.Weather, unsibscribingUrl string) error {
	s.Called(subscriber, weather, unsibscribingUrl)
	return nil
}
