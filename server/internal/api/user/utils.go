package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func GetId(r *http.Request) (result int, err error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, errors.New("id is empty")
	}
	return strconv.Atoi(idStr)
}

func GetUser(r *http.Request) (user *domain.User, err error) {
	err = json.NewDecoder(r.Body).Decode(&user)
	return user, err
}
