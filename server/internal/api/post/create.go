package post

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

func (u *Post) Create(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	post, err := GetPost(request)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	postObject, err := u.service.CreatePost(ctx, post)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, postObject)
}
