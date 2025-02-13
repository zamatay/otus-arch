package user

import (
	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"net/http"
)

func (api *User) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	user, err := GetUser(r)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	id, err := api.service.AddUser(ctx, *user)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, struct {
		id int
	}{id: id})
}
