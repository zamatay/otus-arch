package post

import (
	"context"
	"encoding/json"
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

type PostServiced interface {
	CreatePost(context.Context, *domain.Post) (*domain.Post, error)
	UpdatePost(context.Context, *domain.Post) (bool, error)
	DeletePost(context.Context, int, int) (bool, error)
	GetPost(context.Context, int) (*domain.Post, error)
	FeedPost(context.Context, int, int, int) ([]*domain.Post, error)
}

type Post struct {
	service  PostServiced
	cache    *redis.Cache //FeedPosted
	producer *kafka.Producer
}

func NewPost(service PostServiced, cache *redis.Cache, s *srvApi.Service, producer *kafka.Producer) *Post {
	api := new(Post)
	api.service = service
	api.cache = cache
	api.producer = producer
	api.RegisterHandler(s)
	return api
}

func (u *Post) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/post/create", u.Create)
	r.AddProtectedRoute("/post/update", u.Update)
	r.AddProtectedRoute("/post/delete", u.Delete)
	r.AddProtectedRoute("/post/get", u.Get)
	r.AddProtectedRoute("/post/feed", u.Feed)
}

func GetPost(request *http.Request) (post *domain.Post, err error) {
	err = json.NewDecoder(request.Body).Decode(&post)
	return post, err
}
