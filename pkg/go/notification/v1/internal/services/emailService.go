package services

import (
	"context"
	"fmt"
	"time"

	localConf "github.com/Rabiann/weather-mailer/pkg/go/notification/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/notification/internal/dto"
	external "github.com/Rabiann/weather-mailer/pkg/go/notification/internal/external/mailing"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type (
	MailingService struct {
		Provider             MailingProvider
		ConfirmationTemplate ConfirmationTemplater
		WeatherTemplate      WeatherTemplater
		Config               *localConf.Configuration
	}

	MailingProvider interface {
		SendLetter(options external.MailOptions, ctx context.Context) error
	}

	ConfirmationTemplater interface {
		BuildConfirmationLetter(url string) string
	}

	WeatherTemplater interface {
		BuildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string
	}
)

func NewMailingService(mailProvider MailingProvider, config *localConf.Configuration) (*MailingService, error) {
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
	options := external.MailOptions{
		From:    from,
		To:      to,
		Subject: subject,
		Content: body,
	}

	return s.Provider.SendLetter(options, ctx)
}

func (s *MailingService) SendWeatherReport(subscriber *dto.Subscriber, weather *dto.Weather, unsibscribingUrl string) error {
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
	options := external.MailOptions{
		From:    from,
		To:      to,
		Subject: subject,
		Content: body,
	}

	return s.Provider.SendLetter(options, ctx)
}
