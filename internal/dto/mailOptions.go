package dto

import "github.com/sendgrid/sendgrid-go/helpers/mail"

type MailOptions struct {
	From    mail.Email
	To      mail.Email
	Subject string
	Content string
}
