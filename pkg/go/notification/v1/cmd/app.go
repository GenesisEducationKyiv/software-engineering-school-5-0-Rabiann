package cmd

import (
	"os"

	"github.com/Rabiann/weather-mailer/lib/go/config"
	"github.com/Rabiann/weather-mailer/lib/go/logger"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/notification/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/notification/internal/consumers"
	external "github.com/Rabiann/weather-mailer/pkg/go/notification/internal/external/mailing"
	"github.com/Rabiann/weather-mailer/pkg/go/notification/internal/services"
)

type App struct{}

func (a *App) Run() error {
	configuration := &localConf.Configuration{}
	err := config.LoadEnvironment(configuration)
	if err != nil {
		return err
	}

	os.Mkdir("logs", 0700)
	logger.SetupLogger("logs/app.log")
	mailingProvider := external.NewMailingProvider(configuration)
	emailService, err := services.NewMailingService(mailingProvider, configuration)
	if err != nil {
		return err
	}

	confirmationConsumer, err := consumers.NewConfirmationConsumer(emailService, configuration)
	if err != nil {
		return err
	}

	go confirmationConsumer.StartWorker(nil)

	reportConsumer, err := consumers.NewReportConsumer(emailService, configuration)
	if err != nil {
		return err
	}

	go reportConsumer.StartWorker(nil)

	return nil
}
