package dialogs

import (
	"net/http"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (u *Dialog) Send(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userFrom := domain.GetUserFromContext(ctx)

	d, err := GetDialog(request)
	if err != nil {
		srvApi.SetError(writer, err.Error(), 500)
		return
	}

	if d.Text == "" {
		srvApi.SetError(writer, "Text is empty", 400)
		return
	}

	IsOk, err := u.service.SendDialog(ctx, userFrom.Id, d.ToUserID, d.Text)
	if err != nil {
		return
	}

	srvApi.SetOk(writer, srvApi.OkFalse(IsOk))

}
