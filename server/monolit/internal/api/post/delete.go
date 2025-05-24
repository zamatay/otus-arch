package post

import (
	"log/slog"
	"net/http"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/kafka"
)

func (u *Post) Delete(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	id := request.URL.Query().Get("id")

	userFrom := domain.GetUserFromContext(ctx)
	isOk, err := u.service.DeletePost(ctx, id, userFrom.Id)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	if message, err := kafka.CreateMessage[domain.Post]("", "delete", domain.Post{ID: id, UserID: userFrom.Id}, "posts"); err == nil {
		if err := u.Producer.Produce(message); err != nil {
			slog.Error("Ошибка при отправке kafka", "error", err)
		}
	}

	srvApi.SetOk(writer, srvApi.OkFalse(isOk))
}
