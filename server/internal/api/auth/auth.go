package auth

import (
	"context"
	"encoding/json"
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

var signingKey []byte

type AuthServiced interface {
	Login(context.Context, string, string) (string, error)
	Register(context.Context, domain.RegisterUser) error
}

type Auth struct {
	service AuthServiced
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var au domain.AuthUser
	if err := json.NewDecoder(r.Body).Decode(&au); err != nil {
		srvApi.SetError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := a.service.Login(ctx, au.User, au.Password)
	if err != nil {
		srvApi.SetError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := domain.AuthUserResult{Token: token}
	srvApi.SetOk(w, result)
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var u domain.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		srvApi.SetError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.service.Register(ctx, u); err != nil {
		srvApi.SetError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	srvApi.SetOk(w, srvApi.Ok())
}

func NewAuth(service AuthServiced, secret string) *Auth {
	signingKey = []byte(secret)
	return &Auth{service: service}
}

func (a *Auth) RegisterHandler(r srvApi.AddRouted) {
	r.AddRoute("/auth/login", a.Login)
	r.AddRoute("/auth/register", a.Register)
}
