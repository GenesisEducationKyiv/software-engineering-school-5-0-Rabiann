package producers

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/api/internal/config"
	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
)

type (
	ReportProducer struct {
		producer sarama.SyncProducer
		topic    string
	}
)

func NewReportProducer(config *localConf.Configuration) (*ReportProducer, error) {
	brokers := []string{config.KafkaBrokerUrl}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return nil, err
	}

	return &ReportProducer{
		producer: producer,
		topic:    config.ReportTopic,
	}, nil
}

func (p *EmailProducer) Publish_SendWeatherReport(message dto.ReportMsg) error {
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
