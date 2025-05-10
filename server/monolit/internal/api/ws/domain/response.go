package domain

type OkMessage struct {
	Result bool   `json:"result"`
	Action string `json:"action,omitempty" `
}

func GetOk() OkMessage {
	return OkMessage{Result: true}
}

func (w *OkMessage) SetAction(value string) {
	w.Action = value
}
