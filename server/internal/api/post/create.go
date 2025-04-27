package post

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

const cacheCount = 1000

func (u *Post) Create(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	post, err := GetPost(request)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	postObject, err := u.service.CreatePost(ctx, post)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	/*	//Что бы не задерживать пользователю ответ, кэширование будем проводить конкурентно
		go func(post domain.Post) {
		}(*postObject)*/

	srvApi.SetOk(writer, postObject)
}
