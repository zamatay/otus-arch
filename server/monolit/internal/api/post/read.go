package post

import (
	"net/http"

	srvApi "github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/pkg/api/counter"
)

type ReadResponse struct {
	srvApi.OkResult
	Count int64 `json:"count"`
}

func (u *Post) Read(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	post_id, err := srvApi.GetByName(request, "post_id")
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}
	user_id, err := srvApi.GetByName(request, "user_id")
	if err != nil {
		http.Error(writer, err.Error(), 400)
		return
	}

	count, err := u.service.Read(ctx, post_id, user_id)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	in := &counter.CounterRequest{
		PostId: int32(post_id),
	}

	u.counter.Increase(ctx, in)

	srvApi.SetOk(writer, ReadResponse{Count: count})
}
