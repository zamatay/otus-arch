package api

import (
	"context"
	"net/http"
)

type CounterServiced interface {
	GetPostRead(context.Context, int32) (int, error)
}

type Counter struct {
	service CounterServiced
}

func NewCounter(service CounterServiced, s *Service) *Counter {
	api := new(Counter)
	api.service = service
	api.RegisterHandler(s)
	return api
}

func (u *Counter) RegisterHandler(r AddRouted) {
	r.AddProtectedRoute("/post/Count", u.GetReadPost)
}

func (u *Counter) GetReadPost(writer http.ResponseWriter, request *http.Request) {
	ctx, done := GetContext(request.Context())
	defer done()

	id, err := GetByName(request, "post_id")
	if err != nil {
		SetError(writer, err.Error(), 500)
		return
	}

	postRead, err := u.service.GetPostRead(ctx, int32(id))
	if err != nil {
		SetError(writer, err.Error(), 500)
		return
	}

	SetOk(writer, postRead)
}
