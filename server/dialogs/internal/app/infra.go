package app

import (
	"context"

	"dialogs/internal/api"
	"dialogs/internal/api/grpcclient"
	"dialogs/internal/config"
	"dialogs/internal/repository"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewInfra(ctx context.Context, config *config.Config) (*repository.Repo, *api.Service, error) {
	repo, err := repository.NewRepo(ctx, config.DB)
	if err != nil {
		return nil, nil, err
	}

	service, err := api.New(&config.Http, config.App.Secret)
	if err != nil {
		return nil, nil, err
	}

	client := grpcclient.NewMonolitService(config.GRPC)
	checker := api.NewHealthCheckHandler(client)
	service.AddHandle("/health", checker)
	service.AddHandle("/metrics", promhttp.Handler())

	return repo, service, err

}
