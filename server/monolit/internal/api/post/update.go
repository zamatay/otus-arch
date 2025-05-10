package post

import (
	"net/http"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
)

func (u *Post) Update(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	post, err := GetPost(request)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	isOk, err := u.service.UpdatePost(ctx, post)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, srvApi.OkFalse(isOk))
}
