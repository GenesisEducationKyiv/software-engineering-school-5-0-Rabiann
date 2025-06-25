package external

import (
	"context"

	"github.com/Rabiann/weather-mailer/internal/config"
	"github.com/Rabiann/weather-mailer/internal/dto"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type (
	MailingProvider struct {
		Client *sendgrid.Client
	}
)

func NewMailingProvider(config *config.Configuration) *MailingProvider {
	client := sendgrid.NewSendClient(config.SendgridApiKey)
	return &MailingProvider{client}
}

func (s *MailingProvider) SendLetter(options dto.MailOptions, ctx context.Context) error {
	message := mail.NewSingleEmail(&options.From, options.Subject, &options.To, "", options.Content)
	_, err := s.Client.SendWithContext(ctx, message)
	return err
}
