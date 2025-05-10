package dialogs

import (
	"context"
	"encoding/json"
	"net/http"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

type DialogServiced interface {
	SendDialog(ctx context.Context, fromUserId int, toUserId int, text string) (*domain.Dialog, error)
	ListDialog(ctx context.Context, fromUserId int, toUserId int) ([]*domain.Dialog, error)
}

type Dialog struct {
	service DialogServiced
	cache   *redis.Cache //FeedPosted
}

func NewDialog(service DialogServiced, cache *redis.Cache, s *srvApi.Service) *Dialog {
	d := new(Dialog)
	d.cache = cache
	d.service = service
	d.RegisterHandler(s)
	return &Dialog{service: service, cache: cache}
}

func (u *Dialog) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/dialog/create", u.Send)
	r.AddProtectedRoute("/dialog/list", u.List)
}

func GetDialog(request *http.Request) (dialog *domain.Dialog, err error) {
	err = json.NewDecoder(request.Body).Decode(&dialog)
	return dialog, err
}
