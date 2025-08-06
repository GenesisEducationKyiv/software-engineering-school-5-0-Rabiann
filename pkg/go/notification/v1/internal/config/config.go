package config

type Configuration struct {
	SenderMail        string
	MailTimeout       int
	SendgridApiKey    string
	BaseUrl           string
	KafkaBrokerUrl    string
	ConfirmationTopic string
}
