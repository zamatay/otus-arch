package user

import (
	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"net/http"
)

func (api *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var err error

	user, err := GetUser(r)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	err = api.service.UpdateUser(ctx, *user)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, user)
}
