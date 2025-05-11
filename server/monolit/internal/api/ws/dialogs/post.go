package dialogs

import (
	"context"
	"encoding/json"
	"log/slog"

	domain2 "github.com/zamatay/otus/arch/lesson-1/internal/api/ws/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type DialogMessage struct {
	domain2.OkMessage
	Dialogs domain.Dialog
}

type DialogPostMessage struct {
	domain2.AuthMessage
	Text       string `json:"text"`
	FromUserId int    `json:"from_user_id"`
	ToUserId   int    `json:"to_user_id"`
}

const DialogPostAction = "dialogs/post"

func (w *Dialogs) HandleDialogPost(data string) any {
	mess := DialogPostMessage{}
	err := json.Unmarshal([]byte(data), &mess)
	if err != nil {
		slog.Error("Ошибка при получении сообщения аутентификации", "message", data, "err", err)
		return nil
	}
	dialog, err := w.srv.SendDialog(context.Background(), mess.FromUserId, mess.ToUserId, mess.Text)
	if err != nil {
		return nil
	}
	result := DialogMessage{OkMessage: domain2.OkMessage{Result: true, Action: DialogPostAction}, Dialogs: *dialog}
	return result
}
