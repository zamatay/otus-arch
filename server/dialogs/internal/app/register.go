package app

import (
	"context"

	"dialogs/internal/api"
	"dialogs/internal/api/dialogs"
	"dialogs/internal/repository"
)

func RegisterApi(ctx context.Context, repo *repository.Repo, service *api.Service) {
	dialogs.NewDialog(repo, service)
}
