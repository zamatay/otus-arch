package post

import (
	"log/slog"
	"net/http"
	"strconv"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
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

	if message, err := kafka.CreateMessage[domain.Post](strconv.Itoa(postObject.ID), "create", *postObject); err == nil {
		if err := u.producer.Produce(message); err != nil {
			slog.Error("Ошибка при отправке kafka", "error", err)
		}
	}

	/*	//Что бы не задерживать пользователю ответ, кэширование будем проводить конкурентно
		go func(post domain.Post) {
		}(*postObject)*/

	srvApi.SetOk(writer, postObject)
}
