package dialogs

import (
	"context"
	"encoding/json"
	"net/http"

	srvApi "dialogs/internal/api"
	"dialogs/internal/model"
)

type DialogServiced interface {
	SendDialog(ctx context.Context, fromUserId int, toUserId int, text string) (*model.Dialog, error)
	ListDialog(ctx context.Context, fromUserId int, toUserId int) ([]*model.Dialog, error)
}

type Dialog struct {
	service DialogServiced
}

func NewDialog(service DialogServiced, s *srvApi.Service) *Dialog {
	d := new(Dialog)
	d.service = service
	d.RegisterHandler(s)
	return &Dialog{service: service}
}

func (u *Dialog) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/dialog/create", u.Send)
	r.AddProtectedRoute("/dialog/list", u.List)
}

func GetDialog(request *http.Request) (dialog *model.Dialog, err error) {
	err = json.NewDecoder(request.Body).Decode(&dialog)
	return dialog, err
}
