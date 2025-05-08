package ws

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/gorilla/websocket"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/dialogs"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/posts"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type getterUser interface {
	GetUser() int
}

type setterAction interface {
	SetAction(value string)
}

type basicAction struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func (w WS) Handle(con *websocket.Conn, messageType int, message []byte) error {
	action := basicAction{}
	err := json.Unmarshal(message, &action)
	if err != nil {
		slog.Error("Неизвестный action", "body", string(message))
		return err
	}
	var data any
	switch action.Action {
	case AuthAction:
		data = w.handleAuth(action.Data)
	case dialogs.DialogListAction:
		data = w.dialogs.HandleDialogList(action.Data)
	case dialogs.DialogPostAction:
		data = w.dialogs.HandleDialogPost(action.Data)
	case posts.PostsPostAction:
		data = w.posts.HandlePost(action.Data)
	case posts.PostsListAction:
		data = w.posts.HandleList(action.Data)
	}

	if user, ok := data.(*domain.UserClaims); ok {
		userID := user.Id
		w.activeConnections[userID] = con
	}

	if act, ok := data.(setterAction); ok {
		act.SetAction(action.Action)
	}

	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = con.WriteMessage(messageType, response)
	if err != nil {
		log.Println("Ошибка записи:", err)
		return err
	}

	return nil
}
