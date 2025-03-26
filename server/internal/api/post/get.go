package post

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

func (u *Post) Get(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	id, err := srvApi.GetByName(request, "id")
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	post, err := u.service.GetPost(ctx, id)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, post)
}
