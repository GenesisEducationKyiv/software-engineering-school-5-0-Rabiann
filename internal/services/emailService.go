package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Subscriber struct {
	Recipient string
	Period    string
	City      string
}

type MailingService struct {
	Client               *sendgrid.Client
	ConfirmationTemplate *ConfirmationTemplate
	WeatherTemplate      *WeatherTemplate
	Config               *config.Configuration
}

type MailingServer interface {
	SendConfirmationLetter(string, string) error
	sendLetter(MailOptions, context.Context) error
	SendWeatherReport(*Subscriber, *Weather, string) error
}

type Template struct {
	text string
}

type ConfirmationTemplate struct {
	template *Template
}

type WeatherTemplate struct {
	template *Template
}

type MailOptions struct {
	from    mail.Email
	to      mail.Email
	subject string
	content string
}

func NewTemplate(filepath string) (*Template, error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var template Template

	template.text = string(text)
	return &template, nil
}

func NewConfirmationTemplate(filepath string) (*ConfirmationTemplate, error) {
	template, err := NewTemplate(filepath)
	if err != nil {
		return nil, err
	}

	return &ConfirmationTemplate{template: template}, nil
}

func NewWeatherTemplate(filepath string) (*WeatherTemplate, error) {
	template, err := NewTemplate(filepath)
	if err != nil {
		return nil, err
	}

	return &WeatherTemplate{template: template}, nil
}

func NewMailingService(config *config.Configuration) (*MailingService, error) {
	var ms MailingService
	client := sendgrid.NewSendClient(config.SendgridApiKey)
	ms.Client = client

	confirmationTemplate, err := NewConfirmationTemplate("./templates/confirmationMail.tmpl")
	if err != nil {
		return nil, err
	}

	weatherTemplate, err := NewWeatherTemplate("./templates/weatherMail.tmpl")
	if err != nil {
		return nil, err
	}

	ms.ConfirmationTemplate = confirmationTemplate
	ms.WeatherTemplate = weatherTemplate
	ms.Config = config
	return &ms, nil
}

func (ct *ConfirmationTemplate) buildConfirmationLetter(email string) string {
	return strings.Replace(ct.template.text, "{}", email, 3)
}

func (wt *WeatherTemplate) buildWeatherLetter(city string, temp string, humid string, description string, unsubscribe string) string {
	let := strings.Replace(wt.template.text, "{City}", city, 1)
	let = strings.Replace(let, "{Temperature}", temp, 1)
	let = strings.Replace(let, "{Humidity}", humid, 1)
	let = strings.Replace(let, "{UnsubscribeLink}", unsubscribe, 1)
	let = strings.Replace(let, "{Description}", description, 1)
	return let
}

func (s *MailingService) sendLetter(options MailOptions, ctx context.Context) error {
	message := mail.NewSingleEmail(&options.from, options.subject, &options.to, "", options.content)
	_, err := s.Client.SendWithContext(ctx, message)
	return err
}

func (s *MailingService) SendConfirmationLetter(recipient string, confirmationUrl string) error {
	from := mail.Email{
		Name:    "Confirmator",
		Address: os.Getenv("SENDER_MAIL"),
	}
	to := mail.Email{
		Name:    recipient,
		Address: recipient,
	}

	subject := "Confirm Weather Subscription"
	body := s.ConfirmationTemplate.buildConfirmationLetter(confirmationUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.Config.MailTimeout))
	defer cancel()
	options := MailOptions{
		from:    from,
		to:      to,
		subject: subject,
		content: body,
	}

	return s.sendLetter(options, ctx)
}

func (s *MailingService) SendWeatherReport(subscriber *Subscriber, weather *Weather, unsibscribingUrl string) error {
	from := mail.Email{
		Name:    "Reporter",
		Address: os.Getenv("SENDER_MAIL"),
	}
	to := mail.Email{
		Name:    subscriber.Recipient,
		Address: subscriber.Recipient,
	}

	subject := fmt.Sprintf("%s report for %s", subscriber.Period, subscriber.City)
	body := s.WeatherTemplate.buildWeatherLetter(subscriber.City, fmt.Sprintf("%.1f", weather.Temperature), fmt.Sprintf("%.1f", weather.Humidity), weather.Description, unsibscribingUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	options := MailOptions{
		from:    from,
		to:      to,
		subject: subject,
		content: body,
	}

	return s.sendLetter(options, ctx)
}
