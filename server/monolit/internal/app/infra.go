package app

import (
	"context"

	"github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/grpcclient/counter"
	"github.com/zamatay/otus/arch/lesson-1/internal/config"
	"github.com/zamatay/otus/arch/lesson-1/internal/grpcserver"
	"github.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func NewInfra(ctx context.Context, config *config.Config) (*repository.Repo, *redis.Cache, *api.Service, *kafka.Producer, error) {
	repo, err := repository.NewRepo(ctx, config.DB["read"], config.DB["write"], config.DB["shard"])
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cache, err := redis.NewCache(ctx, config.Cache)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	kafkaProducer, err := kafka.NewProducer(config.Kafka)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	service, err := api.New(&config.Http, config.App.Secret)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	service.Grpc, err = grpcserver.NewGRPCServer(&config.GRPC)

	service.CounterSrv = counter.NewCounterService(config.GRPCCounter)

	return repo, cache, service, kafkaProducer, err

}
