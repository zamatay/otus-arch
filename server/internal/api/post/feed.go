package post

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
)

func (u *Post) Feed(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	offset, _ := srvApi.GetByName(request, "offset")
	limit, err := srvApi.GetByName(request, "limit")
	if err != nil {
		limit = 100
	}

	posts, err := u.service.FeedPost(ctx, offset, limit)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, posts)
}
