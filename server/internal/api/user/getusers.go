package user

import (
	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"net/http"
)

func (api *User) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	srvApi.SetOk(w, api.service.GetUsers(ctx))
}

func (api *User) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var id int
	var err error
	if id, err = GetId(r); err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, api.service.GetUser(ctx, id))

}
