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
}

type User struct {
	service UserServiced
}

func NewUser(service UserServiced) *User {
	api := new(User)
	api.service = service

	return api
}

func (u *User) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/user/get_list", u.GetUsers)
	r.AddProtectedRoute("/user/add", u.AddUser)
	r.AddProtectedRoute("/user/update", u.UpdateUser)
	r.AddProtectedRoute("/user/get", u.GetUser)
	r.AddProtectedRoute("/user/remove", u.Remove)
}

func run(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Ошибка при попытки запустить http сервер", err)
	}
}

func (api *User) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	srvApi.SetOk(w, api.service.GetUsers(ctx))
}

func (api *User) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var id int
	var err error
	if id, err = GetId(r); err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, api.service.GetUser(ctx, id))

}

func (api *User) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	user, err := GetUser(r)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	id, err := api.service.AddUser(ctx, *user)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, struct {
		id int
	}{id: id})
}

func (api *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var err error

	user, err := GetUser(r)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	err = api.service.UpdateUser(ctx, *user)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, user)
}

func (api *User) Remove(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var id int
	var err error
	if id, err = GetId(request); err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	if err = api.service.Remove(ctx, id); err != nil {
		srvApi.SetError(writer, err.Error(), http.StatusInternalServerError)
	}

	srvApi.SetOk(writer, srvApi.Ok())
}
