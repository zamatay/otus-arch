package app

import (
	"counter/internal/api"
	"counter/internal/config"
	"counter/internal/grpcserver"
	"counter/internal/repository"
)

func NewService(config *config.Config, repo *repository.Repo) (*api.Service, *grpcserver.Service, error) {
	service, err := api.New(&config.Http, config.App.Secret)
	if err != nil {
		return nil, nil, err
	}

	api.NewCounter(repo, service)

	service.AddHandle("/health", api.NewHealthCheckHandler())

	grpc, err := grpcserver.NewGRPCServer(&config.GRPCServer, repo)

	return service, grpc, nil
}
