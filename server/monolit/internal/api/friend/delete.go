package friend

import (
	"net/http"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	userSrv "githib.com/zamatay/otus/arch/lesson-1/internal/api/user"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (u *Friend) Delete(w http.ResponseWriter, r *http.Request) {
	user := domain.GetUserFromContext(r.Context())
	if user == nil {
		api.SetError(w, "User not found", 404)
		return
	}

	id, err := userSrv.GetId(r)
	if err != nil {
		api.SetError(w, "User not found", 404)
		return
	}

	err = u.service.DeleteFriends(r.Context(), user.Id, id)
	if err != nil {
		api.SetError(w, err.Error(), 500)
		return
	}

	api.SetOk(w, api.Ok())
}
