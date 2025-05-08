package dialogs

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api/utils"
	domain2 "githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type DialogsMessage struct {
	domain2.OkMessage
	Dialogs []*domain.Dialog
}
type DialogListMessage struct {
	domain2.AuthMessage
	UserId string `json:"user_id"`
}

const DialogListAction = "dialogs/list"

func (w Dialogs) HandleDialogList(data string) any {
	mess := DialogListMessage{}
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
	toUserId, err := strconv.Atoi(mess.UserId)
	dialog, err := w.srv.ListDialog(context.Background(), user.Id, toUserId)
	if err != nil {
		return nil
	}

	dm := DialogsMessage{
		OkMessage: domain2.OkMessage{Result: true},
		Dialogs:   dialog,
	}

	dm.SetAction(DialogListAction)

	return dm
}
