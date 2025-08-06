package producers

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/api/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
)

type (
	EmailProducer struct {
		producer sarama.SyncProducer
		topic    string
	}
)

func NewEmailProducer(config *localConf.Configuration) (*EmailProducer, error) {
	brokers := []string{config.KafkaBrokerUrl}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return nil, err
	}

	return &EmailProducer{
		producer: producer,
		topic:    config.ConfirmationTopic,
	}, nil
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func (p *EmailProducer) Publish_SendConfirmationLetter(message dto.SubscriptionConfirmationMsg) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(bytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Event is stored in topic (%s)/partition(%d)/offset(%d)\n", p.topic, partition, offset)
	return nil
}
