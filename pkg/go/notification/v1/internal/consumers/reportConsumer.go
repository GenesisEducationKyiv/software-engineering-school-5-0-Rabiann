package consumers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/notification/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/notification/internal/dto"
)

type (
	ReportConsumer struct {
		consumer     sarama.Consumer
		emailService EmailService_Report
		topic        string
	}

	EmailService_Report interface {
		SendWeatherReport(subscriber *dto.Subscriber, weather *dto.Weather, unsibscribingUrl string) error
	}
)

func NewReportConsumer(emailService EmailService_Report, config *localConf.Configuration) (*ReportConsumer, error) {
	brokers := []string{config.KafkaBrokerUrl}
	consumer, err := ConnectConsumer(brokers)
	if err != nil {
		return nil, err
	}

	return &ReportConsumer{
		consumer:     consumer,
		emailService: emailService,
		topic:        config.ConfirmationTopic,
	}, nil
}

func Parse_Report(data []byte) (dto.ReportMsg, error) {
	var confirmation dto.ReportMsg
	err := json.Unmarshal(data, &confirmation)

	return confirmation, err
}

func (c *ReportConsumer) StartWorker(errCh chan error) error {
	consumer := c.consumer
	worker, err := consumer.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		errCh <- err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-worker.Errors():
				log.Fatalf("%v", err)
			case msg := <-worker.Messages():
				log.Print(msg)
				data := msg.Value
				report, err := Parse_Report(data)
				if err != nil {
					log.Fatal(err)
				}
				c.emailService.SendWeatherReport(&report.Subscriber, &report.Weather, report.UnsubscriptionUrl)
			case <-sigchan:
				log.Fatal("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Fatal("Stopping consumer")
	return consumer.Close()
}
