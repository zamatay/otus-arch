package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"githib.com/zamatay/otus/arch/lesson-1/internal/app"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type UserServiced interface {
	GetUsers(context.Context) []domain.User
	GetUser(context.Context, int) *domain.User
	AddUser(context.Context, domain.User) (int, error)
	UpdateUser(context.Context, int, domain.User) error
	Remove(context.Context, int) error
}

type User struct {
	service UserServiced
}

const _timeout time.Duration = 2 * time.Second

func NewUser(service UserServiced, config app.HttpConfig) *User {
	api := new(User)
	http.HandleFunc("/user/get_list", api.GetUsers)
	http.HandleFunc("/user/add", api.AddUser)
	http.HandleFunc("/user/update", api.UpdateUser)
	http.HandleFunc("/user/get", api.GetUser)
	http.HandleFunc("/user/remove", api.Remove)
	api.service = service

	go run(config.Host, config.Port)

	return api
}

func run(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Ошибка при попытки запустить http сервер", err)
	}
}

func (api *User) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, done := context.WithTimeout(context.Background(), _timeout)
	defer done()
	err := json.NewEncoder(w).Encode(api.service.GetUsers(ctx))
	if err != nil {
		setError(w, err.Error(), 500)
	}
	w.WriteHeader(http.StatusOK)
}

func (api *User) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := context.WithTimeout(context.Background(), _timeout)
	defer done()

	var id int
	var err error
	if id, err = GetId(r); err != nil {
		setError(w, err.Error(), 500)
		return
	}

	setOk(w, api.service.GetUser(ctx, id))

}

func (api *User) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := context.WithTimeout(context.Background(), _timeout)
	defer done()

	user, err := GetUser(r)
	if err != nil {
		setError(w, err.Error(), 500)
	}

	id, err := api.service.AddUser(ctx, *user)
	if err != nil {
		setError(w, err.Error(), 500)
	}

	setOk(w, struct {
		id int
	}{id: id})
}

func (api *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := context.WithTimeout(context.Background(), _timeout)
	defer done()

	var id int
	var err error
	if id, err = GetId(r); err != nil {
		setError(w, err.Error(), 500)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		setError(w, err.Error(), 500)
		return
	}

	err = api.service.UpdateUser(ctx, id, *user)
	if err != nil {
		setError(w, err.Error(), 500)
		return
	}

	setOk(w, user)
}

func (api *User) Remove(writer http.ResponseWriter, request *http.Request) {
	ctx, done := context.WithTimeout(context.Background(), _timeout)
	defer done()

	var id int
	var err error
	if id, err = GetId(request); err != nil {
		setError(writer, err.Error(), 500)
		return
	}

	setOk(writer, api.service.Remove(ctx, id))
}
