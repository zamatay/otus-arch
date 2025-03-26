package post

import (
	"context"
	"encoding/json"
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type PostServiced interface {
	CreatePost(context.Context, *domain.Post) (*domain.Post, error)
	UpdatePost(context.Context, *domain.Post) (bool, error)
	DeletePost(context.Context, int) (bool, error)
	GetPost(context.Context, int) (*domain.Post, error)
	FeedPost(context.Context, int, int) ([]*domain.Post, error)
}

type Post struct {
	service PostServiced
}

func NewPost(service PostServiced, s *srvApi.Service) *Post {
	api := new(Post)
	api.service = service
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
