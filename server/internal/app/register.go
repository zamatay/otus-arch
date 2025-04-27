package app

import (
	"context"
	"log"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/auth"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/dialogs"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/friend"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/post"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/user"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/ws"
	"githib.com/zamatay/otus/arch/lesson-1/internal/config"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func RegisterApi(ctx context.Context, repo *repository.Repo, service *api.Service, cache *redis.Cache, secret string, config *config.Config) {
	user.NewUser(repo, service)
	auth.NewAuth(repo, service, secret)
	friend.NewFriend(repo, service)
	post.NewPost(repo, cache, service)
	dialogs.NewDialog(repo, cache, service)
	_, err := ws.NewWS(ctx, repo, service, config)
	if err != nil {
		log.Fatal(err)
	}
}
