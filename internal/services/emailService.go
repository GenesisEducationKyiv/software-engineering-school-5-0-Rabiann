package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/dto"
	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type (
	MailingService struct {
		Provider             MailingProvider
		ConfirmationTemplate ConfirmationTemplater
		WeatherTemplate      WeatherTemplater
		Config               *config.Configuration
	}

	MailingProvider interface {
		SendLetter(options dto.MailOptions, ctx context.Context) error
	}

	ConfirmationTemplater interface {
		BuildConfirmationLetter(url string) string
	}

	WeatherTemplater interface {
		BuildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string
	}
)

func NewMailingService(mailProvider MailingProvider, config *config.Configuration) (*MailingService, error) {
	var ms MailingService
	ms.Provider = mailProvider

	confirmationTemplate, err := NewConfirmationTemplate("templates/confirmationMail.tmpl")
	if err != nil {
		return nil, err
	}

	weatherTemplate, err := NewWeatherTemplate("templates/weatherMail.tmpl")
	if err != nil {
		return nil, err
	}

	ms.ConfirmationTemplate = confirmationTemplate
	ms.WeatherTemplate = weatherTemplate
	ms.Config = config
	return &ms, nil
}

func (s *MailingService) SendConfirmationLetter(recipient string, confirmationUrl string) error {
	from := mail.Email{
		Name:    "Confirmator",
		Address: s.Config.SenderMail,
	}
	to := mail.Email{
		Name:    recipient,
		Address: recipient,
	}

	subject := "Confirm Weather Subscription"
	body := s.ConfirmationTemplate.BuildConfirmationLetter(confirmationUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.Config.MailTimeout))
	defer cancel()
	options := dto.MailOptions{
		From:    from,
		To:      to,
		Subject: subject,
		Content: body,
	}

	return s.Provider.SendLetter(options, ctx)
}

func (s *MailingService) SendWeatherReport(subscriber *models.Subscriber, weather *models.Weather, unsibscribingUrl string) error {
	from := mail.Email{
		Name:    "Reporter",
		Address: s.Config.SenderMail,
	}
	to := mail.Email{
		Name:    subscriber.Recipient,
		Address: subscriber.Recipient,
	}

	subject := fmt.Sprintf("%s report for %s", subscriber.Period, subscriber.City)
	body := s.WeatherTemplate.BuildWeatherLetter(subscriber.City, fmt.Sprintf("%.1f", weather.Temperature), fmt.Sprintf("%.1f", weather.Humidity), weather.Description, unsibscribingUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	options := dto.MailOptions{
		From:    from,
		To:      to,
		Subject: subject,
		Content: body,
	}

	return s.Provider.SendLetter(options, ctx)
}
