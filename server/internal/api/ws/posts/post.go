package posts

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	domainmes "githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

const PostsPostAction = "posts/create"

func (w *Posts) HandlePost(data string) any {
	mess := domainmes.PostPostMessage{}
	err := json.Unmarshal([]byte(data), &mess)
	if err != nil {
		slog.Error("Ошибка при получении сообщения аутентификации", "message", data, "err", err)
		return nil
	}
	postInsert := createMessage(mess)
	post, err := w.srv.CreatePost(context.Background(), &postInsert)
	if err != nil {
		slog.Error("Ошибка при создании post", "message", data, "err", err)
		return nil
	}
	result := domainmes.PostMessage{OkMessage: domainmes.OkMessage{Result: true, Action: PostsPostAction}, Post: *post}
	return result
}

func createMessage(mess domainmes.PostPostMessage) domain.Post {
	result := domain.Post{}
	result.UserID = mess.UserID
	result.Text = mess.Text
	result.CreatedAt = time.Now()
	return result
}
