package app

import (
	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/auth"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/dialogs"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/friend"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/post"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/user"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func RegisterApi(repo *repository.Repo, service *api.Service, cache *redis.Cache, producer *kafka.Producer, secret string) {
	user.NewUser(repo, service)
	auth.NewAuth(repo, service, secret)
	friend.NewFriend(repo, service)
	post.NewPost(repo, cache, service, producer)
	dialogs.NewDialog(repo, cache, service)
}
