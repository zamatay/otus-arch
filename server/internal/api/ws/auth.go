package ws

import (
	"encoding/json"
	"log/slog"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api/utils"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/domain"
)

const AuthAction = "auth"

func (w *WS) AddConnection() {

}

func (w WS) handleAuth(data string) any {
	authMessage := domain.AuthMessage{}
	err := json.Unmarshal([]byte(data), &authMessage)
	if err != nil {
		slog.Error("Ошибка при получении сообщения аутентификации", "message", data, "err", err)
		return nil
	}
	if authMessage.Token == "" {
		slog.Error("Токен пустой", "message", data)
		return nil
	}
	user, err := utils.UserByToken(authMessage.Token)
	if err != nil {
		slog.Error("Ошибка при получении пользователя", "token", authMessage.Token, "err", err)
		return nil
	}
	return user
}
