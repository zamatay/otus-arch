package post

import (
	"net/http"

	"github.com/hashicorp/go-uuid"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

const cacheCount = 1000

func (u *Post) Create(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	post, err := GetPost(request)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	if post.ID == "" {
		post.ID, _ = uuid.GenerateUUID()
	}
	postObject, err := u.service.CreatePost(ctx, post)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, postObject)
}
