package app

import (
	"context"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func NewInfra(ctx context.Context, config *Config) (*repository.Repo, *redis.Cache, *kafka.Producer, *api.Service, error) {
	repo, err := repository.NewRepo(ctx, config.DB["read"], config.DB["write"], config.DB["shard"])
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cache, err := redis.NewCache(ctx, config.Cache)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	producer, err := kafka.NewProducer(config.Kafka)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	service, err := api.New(&config.Http, config.App.Secret)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return repo, cache, producer, service, err

}
