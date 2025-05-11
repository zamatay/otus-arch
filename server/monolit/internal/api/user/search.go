package user

import (
	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"net/http"
)

func (api *User) SearchUser(w http.ResponseWriter, r *http.Request) {
	ctx, done := srvApi.GetContext(nil)
	defer done()

	firstName, err := GetByName(r, "first_name")
	if err != nil {
		srvApi.SetError(w, err.Error(), http.StatusBadRequest)
		return
	}

	lastName, err := GetByName(r, "last_name")
	if err != nil {
		srvApi.SetError(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := api.service.SearchUser(ctx, firstName, lastName)
	if err != nil {
		srvApi.SetError(w, err.Error(), 500)
		return
	}

	srvApi.SetOk(w, users)
}
