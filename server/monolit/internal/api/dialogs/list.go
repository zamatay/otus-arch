package dialogs

import (
	"net/http"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (u *Dialog) List(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userFrom := domain.GetUserFromContext(ctx)
	userTo, err := api.GetByName(request, "user_id")
	if err != nil {
		api.SetError(writer, err.Error(), 400)
		return
	}

	dialogs, err := u.service.ListDialog(ctx, userFrom.Id, userTo)
	if err != nil {
		api.SetError(writer, err.Error(), 500)
		return
	}

	api.SetOk(writer, dialogs)
}
