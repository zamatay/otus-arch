package user

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

func (api *User) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext(nil)
	defer done()

	srvApi.SetOk(w, api.service.GetUsers(ctx))
}

func (api *User) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext(nil)
	defer done()

	var id int
	var err error
	if id, err = GetId(r); err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, api.service.GetUser(ctx, id))

}
