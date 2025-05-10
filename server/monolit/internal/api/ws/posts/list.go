package posts

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/zamatay/otus/arch/lesson-1/internal/api/utils"
	domainmes "github.com/zamatay/otus/arch/lesson-1/internal/api/ws/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type ListMessage struct {
	domainmes.OkMessage
	Post []*domain.Post
}

type PostListMessage struct {
	domainmes.AuthMessage
}

const PostsListAction = "posts/list"

func (w *Posts) HandleList(data string) any {
	mess := domainmes.AuthMessage{}
	err := json.Unmarshal([]byte(data), &mess)
	if err != nil {
		slog.Error("Ошибка при получении сообщения аутентификации", "message", data, "err", err)
		return nil
	}
	if mess.Token == "" {
		slog.Error("Токен пустой", "message", data)
		return nil
	}
	user, err := utils.UserByToken(mess.Token)
	if err != nil {
		slog.Error("Ошибка при получении пользователя", "token", mess.Token, "err", err)
		return nil
	}

	list, err := w.srv.FeedPost(context.Background(), 0, 100, user.Id)
	if err != nil {
		slog.Error("Ошибка при создании post", "message", data, "err", err)
		return nil
	}
	return ListMessage{OkMessage: domainmes.OkMessage{Result: true, Action: PostsPostAction}, Post: list}
}
