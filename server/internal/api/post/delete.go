package post

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (u *Post) Delete(writer http.ResponseWriter, request *http.Request) {
	ctx, done := srvApi.GetContext(request.Context())
	defer done()

	id, err := srvApi.GetByName(request, "id")
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}
	userFrom := domain.GetUserFromContext(ctx)
	isOk, err := u.service.DeletePost(ctx, id, userFrom.Id)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	srvApi.SetOk(writer, srvApi.OkFalse(isOk))
}
