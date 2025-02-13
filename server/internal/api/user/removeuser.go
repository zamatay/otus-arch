package user

import (
	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"net/http"
)

func (api *User) Remove(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext()
	defer done()

	var id int
	var err error
	if id, err = GetId(request); err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	if err = api.service.Remove(ctx, id); err != nil {
		srvApi.SetError(writer, err.Error(), http.StatusInternalServerError)
	}

	srvApi.SetOk(writer, srvApi.Ok())
}
