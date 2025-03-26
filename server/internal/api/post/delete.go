package post

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

func (u *Post) Delete(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	id, err := srvApi.GetByName(request, "id")
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	isOk, err := u.service.DeletePost(ctx, id)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, srvApi.OkFalse(isOk))
}
