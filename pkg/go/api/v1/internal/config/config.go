package config

type Configuration struct {
	ConfirmationTopic string
	ReportTopic       string
	KafkaBrokerUrl    string
	Port              string
	BaseUrl           string
}
