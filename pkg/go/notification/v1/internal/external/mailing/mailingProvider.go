package external

import (
	"context"

	localConf "github.com/Rabiann/weather-mailer/pkg/go/notification/internal/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type (
	MailingProvider struct {
		Client *sendgrid.Client
	}
)

func NewMailingProvider(config *localConf.Configuration) *MailingProvider {
	client := sendgrid.NewSendClient(config.SendgridApiKey)
	return &MailingProvider{client}
}

func (s *MailingProvider) SendLetter(options MailOptions, ctx context.Context) error {
	message := mail.NewSingleEmail(&options.From, options.Subject, &options.To, "", options.Content)
	_, err := s.Client.SendWithContext(ctx, message)
	return err
}
