package friend

import (
	"context"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
)

type FriendServiced interface {
	SetFriends(context.Context, int, int) (bool, error)
	DeleteFriends(context.Context, int, int) error
}

type Friend struct {
	service FriendServiced
}

func NewFriend(service FriendServiced, s *srvApi.Service) *Friend {
	api := new(Friend)
	api.service = service
	api.RegisterHandler(s)
	return api
}

func (u *Friend) RegisterHandler(r srvApi.AddRouted) {
	r.AddProtectedRoute("/friend/set", u.Set)
	r.AddProtectedRoute("/friend/delete", u.Delete)
}
