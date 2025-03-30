package main

import (
	"context"
	"log"

	"github.com/IBM/sarama"

	"githib.com/zamatay/otus/arch/lesson-1/internal/app"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cache, err := redis.NewCache(context.Background(), cfg.Cache)

	newConsumer, err := kafka.NewConsumer(&cfg.Kafka, cache)
	if err != nil {
		log.Fatal(err)
	}
	defer newConsumer.Close()
	newConsumer.Process()
}
