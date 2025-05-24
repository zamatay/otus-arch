package post

import (
	"log/slog"
	"net/http"

	"github.com/hashicorp/go-uuid"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/kafka"
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

	if post.ID == "" {
		post.ID, _ = uuid.GenerateUUID()
	}
	postObject, err := u.service.CreatePost(ctx, post)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	if message, err := kafka.CreateMessage[domain.Post](postObject.ID, "posts/create", *postObject, "posts"); err == nil {
		if err := u.Producer.Produce(message); err != nil {
			slog.Error("Ошибка при отправке kafka", "error", err)
		}
	}

	srvApi.SetOk(writer, postObject)
}
