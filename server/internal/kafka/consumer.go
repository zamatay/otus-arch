package kafka

import (
	"context"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
)

type Consumer struct {
	config         *sarama.Config
	instance       sarama.Consumer
	partitionPosts sarama.PartitionConsumer
	messageChanel  chan<- *sarama.ConsumerMessage
}

func NewConsumer(cfg *Config, messageChan chan<- *sarama.ConsumerMessage) (*Consumer, error) {
	kafka := new(Consumer)

	kafka.config = sarama.NewConfig()
	kafka.config.Consumer.Return.Errors = true

	var err error
	if kafka.instance, err = sarama.NewConsumer([]string{cfg.GetAddress()}, kafka.config); err != nil {
		return nil, err
	}

	kafka.partitionPosts, err = kafka.instance.ConsumePartition("posts", 0, sarama.OffsetNewest)
	kafka.messageChanel = messageChan

	return kafka, nil
}

func (k *Consumer) Close() error {
	err := k.partitionPosts.Close()
	err = k.instance.Close()
	return err
}

func (k *Consumer) Process(ctx context.Context) {

	for {
		select {
		case msg := <-k.partitionPosts.Messages():
			log.Printf("Received message: %s\n", string(msg.Value))

			k.messageChanel <- msg

		case err := <-k.partitionPosts.Errors():
			log.Println("Consumer error:", err)

		case <-ctx.Done():
			slog.Info("Приложение завершило работу")
			return
		}
	}

}
