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
	ConfirmationConsumer struct {
		consumer     sarama.Consumer
		emailService EmailService_Confirmation
		topic        string
	}

	EmailService_Confirmation interface {
		SendConfirmationLetter(recipient string, confirmationUrl string) error
	}
)

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}

func NewConfirmationConsumer(emailService EmailService_Confirmation, config *localConf.Configuration) (*ConfirmationConsumer, error) {
	brokers := []string{config.KafkaBrokerUrl}
	consumer, err := ConnectConsumer(brokers)
	if err != nil {
		return nil, err
	}

	return &ConfirmationConsumer{
		consumer:     consumer,
		emailService: emailService,
		topic:        config.ConfirmationTopic,
	}, nil
}

func Parse_Confirmation(data []byte) (dto.SubscriptionConfirmationMsg, error) {
	var confirmation dto.SubscriptionConfirmationMsg
	err := json.Unmarshal(data, &confirmation)

	return confirmation, err
}

func (c *ConfirmationConsumer) StartWorker(errCh chan error) error {
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
				confirmation, err := Parse_Confirmation(data)
				if err != nil {
					log.Fatal(err)
				}
				c.emailService.SendConfirmationLetter(confirmation.Email, confirmation.Url)
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
