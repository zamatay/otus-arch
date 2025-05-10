package domain

import (
	"encoding/json"
	"strconv"

	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type PostMessage struct {
	OkMessage
	Post domain.Post
}

type PostPostMessage struct {
	AuthMessage
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

func (receiver *PostPostMessage) UnmarshalJSON(data []byte) error {
	temp := struct {
		Text   string `json:"text"`
		UserID string `json:"user_id"`
	}{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	receiver.Text = temp.Text
	atoi, err := strconv.Atoi(temp.UserID)
	if err != nil {
		return err
	}
	receiver.UserID = atoi
	return nil
}
