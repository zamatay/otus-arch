package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
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

func run(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Ошибка при попытки запустить http сервер", err)
	}
}
