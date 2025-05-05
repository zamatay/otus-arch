package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer struct {
	config   *sarama.Config
	instance sarama.SyncProducer
}

func NewProducer(cfg Config) (*Producer, error) {
	producer := new(Producer)
	producer.config = sarama.NewConfig()
	producer.config.Producer.RequiredAcks = sarama.WaitForAll
	producer.config.Producer.Return.Successes = true

	var err error
	if producer.instance, err = sarama.NewSyncProducer([]string{cfg.GetAddress()}, producer.config); err != nil {
		return nil, err
	}

	return producer, nil
}

func (p *Producer) Close() error {
	return p.instance.Close()
}

func (p *Producer) Produce(msg *sarama.ProducerMessage) error {
	_, _, err := p.instance.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func CreateMessage[T any](key string, messageType string, value T, topic string) (*sarama.ProducerMessage, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("message-type"),
				Value: []byte(messageType),
			},
		},
		Topic: topic,
		Value: sarama.ByteEncoder(jsonData),
		Key:   sarama.StringEncoder("key"),
	}, nil
}
