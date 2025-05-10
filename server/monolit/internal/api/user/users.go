package user

import (
	"context"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type UserServiced interface {
	GetUsers(context.Context) []domain.User
	GetUser(context.Context, int) *domain.User
	AddUser(context.Context, domain.User) (int, error)
	UpdateUser(context.Context, domain.User) error
	Remove(context.Context, int) error
	SearchUser(context.Context, string, string) ([]domain.User, error)
}

type User struct {
	service UserServiced
}

func NewUser(service UserServiced, s *srvApi.Service) *User {
	api := new(User)
	api.service = service
	api.RegisterHandler(s)
	return api
}

func (u *User) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/user/get_list", u.GetUsers)
	r.AddProtectedRoute("/user/add", u.AddUser)
	r.AddProtectedRoute("/user/update", u.UpdateUser)
	r.AddProtectedRoute("/user/get", u.GetUser)
	r.AddProtectedRoute("/user/remove", u.Remove)
	r.AddProtectedRoute("/user/search", u.SearchUser)
}
