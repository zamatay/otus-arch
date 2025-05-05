package kafka

import (
	"github.com/IBM/sarama"
)

func GetChannel() chan *sarama.ConsumerMessage {
	return make(chan *sarama.ConsumerMessage, 100)
}
