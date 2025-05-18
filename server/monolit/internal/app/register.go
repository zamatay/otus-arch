package app

import (
	"context"
	"log"

	"github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/auth"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/dialogs"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/friend"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/post"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/user"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/ws"
	"github.com/zamatay/otus/arch/lesson-1/internal/config"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func RegisterApi(ctx context.Context, repo *repository.Repo, service *api.Service, cache *redis.Cache, secret string, config *config.Config) {
	user.NewUser(repo, service)
	auth.NewAuth(repo, service, secret)
	friend.NewFriend(repo, service)
	post.NewPost(cache, cache, service, service.CounterSrv.Client)
	dialogs.NewDialog(repo, cache, service)
	_, err := ws.NewWS(ctx, repo, service, config)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Grpc.Register(ctx, service)
	if err != nil {
		log.Fatal(err)
	}
}
