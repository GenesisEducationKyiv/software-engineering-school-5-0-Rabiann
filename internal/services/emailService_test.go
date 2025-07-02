package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/dto"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/Rabiann/weather-mailer/internal/services"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	MockMailingProvider struct {
		mock.Mock
	}

	MockConfirmationTemplate struct {
		mock.Mock
	}

	MockWeatherTemplate struct {
		mock.Mock
	}
)

func (m *MockMailingProvider) SendLetter(options dto.MailOptions, ctx context.Context) error {
	res := m.Called(options, ctx)
	return res.Error(0)
}

func (ct *MockConfirmationTemplate) BuildConfirmationLetter(email string) string {
	ct.Called(email)
	return "BODY"
}

func (ct *MockWeatherTemplate) BuildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string {
	ct.Called(city, temp, humid, description, unsubscribe)
	return "WEATHER"
}

func TestSendConfirmationLetter(t *testing.T) {
	recipient := "RECIPIENT"
	service, providerMock := BuildMocks(recipient)
	err := service.SendConfirmationLetter(recipient, "CONFIRMATION URL")
	require.Nil(t, err)

	require.Equal(t, len(providerMock.Calls), 1)
}

func TestSendWeatherReport(t *testing.T) {
	recipient := "RECIPIENT"
	service, providerMock := BuildMocks(recipient)
	subscriber := &models.Subscriber{
		Recipient: recipient,
		Period:    "string",
		City:      "string",
	}
	weather := &models.Weather{
		Temperature: 1.0,
		Humidity:    1.0,
		Description: "string",
	}

	err := service.SendWeatherReport(subscriber, weather, "UNSUBSCRIBING URL")
	require.Nil(t, err)

	require.Equal(t, len(providerMock.Calls), 1)
	options := dto.MailOptions{
		From: mail.Email{
			Name:    "Reporter",
			Address: "mail@mail.com",
		},
		To: mail.Email{
			Name:    "RECIPIENT",
			Address: "RECIPIENT",
		},
		Subject: "string report for string",
		Content: "WEATHER",
	}
	providerMock.AssertCalled(t, "SendLetter", options, mock.Anything)
}

func BuildConfiguration() *config.Configuration {
	return &config.Configuration{
		SendgridApiKey: "testkey",
		SenderMail:     "mail@mail.com",
	}
}

func BuildMocks(recipient string) (
	services.MailingService,
	*MockMailingProvider,
) {
	subject := "Confirm Weather Subscription"
	providerMock := new(MockMailingProvider)
	mockConfTemp := new(MockConfirmationTemplate)
	mockWeatherTemp := new(MockWeatherTemplate)

	configuration := BuildConfiguration()
	from := mail.Email{
		Name:    "Confirmator",
		Address: configuration.SenderMail,
	}
	to := mail.Email{
		Name:    recipient,
		Address: recipient,
	}
	options1 := dto.MailOptions{
		From:    from,
		To:      to,
		Subject: subject,
		Content: "BODY",
	}

	from2 := mail.Email{
		Name:    "Reporter",
		Address: configuration.SenderMail,
	}
	to2 := mail.Email{
		Name:    recipient,
		Address: recipient,
	}
	options2 := dto.MailOptions{
		From:    from2,
		To:      to2,
		Subject: fmt.Sprintf("%s report for %s", "string", "string"),
		Content: "WEATHER",
	}

	providerMock.On("SendLetter", options1, mock.Anything).Return(nil)
	providerMock.On("SendLetter", options2, mock.Anything).Return(nil)
	mockConfTemp.On("BuildConfirmationLetter", "CONFIRMATION URL").Return("BODY")
	mockWeatherTemp.On("BuildWeatherLetter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("BODY")
	service := services.MailingService{
		Provider:             providerMock,
		ConfirmationTemplate: mockConfTemp,
		WeatherTemplate:      mockWeatherTemp,
		Config:               configuration,
	}

	return service, providerMock
}
