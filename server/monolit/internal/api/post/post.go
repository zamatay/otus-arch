package post

import (
	"context"
	"encoding/json"
	"net/http"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
	"github.com/zamatay/otus/arch/lesson-1/pkg/api/counter"
)

type PostServiced interface {
	CreatePost(context.Context, *domain.Post) (*domain.Post, error)
	UpdatePost(context.Context, *domain.Post) (bool, error)
	DeletePost(context.Context, string, int) (bool, error)
	GetPost(context.Context, string) (*domain.Post, error)
	FeedPost(context.Context, int, int, int) ([]*domain.Post, error)
	Read(context.Context, int, int) (int64, error)
}

type Post struct {
	service PostServiced
	counter counter.CounterClient
	cache   *redis.Cache //FeedPosted
}

func NewPost(service PostServiced, cache *redis.Cache, s *srvApi.Service, srv counter.CounterClient) *Post {
	api := new(Post)
	api.service = service
	api.cache = cache
	api.RegisterHandler(s)
	api.counter = srv
	return api
}

func (u *Post) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/post/create", u.Create)
	r.AddProtectedRoute("/post/update", u.Update)
	r.AddProtectedRoute("/post/delete", u.Delete)
	r.AddProtectedRoute("/post/get", u.Get)
	r.AddProtectedRoute("/post/feed", u.Feed)
	r.AddProtectedRoute("/post/read", u.Read)
}

func GetPost(request *http.Request) (post *domain.Post, err error) {
	err = json.NewDecoder(request.Body).Decode(&post)
	return post, err
}
