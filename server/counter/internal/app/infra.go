package app

import (
	"context"

	"counter/internal/config"
	"counter/internal/repository"
)

func NewInfra(ctx context.Context, config *config.Config) *repository.Repo {
	repo := repository.NewRepo(ctx, &config.Cache)

	return repo

}
