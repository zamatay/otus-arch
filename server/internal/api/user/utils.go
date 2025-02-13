package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

var IsEmpty = errors.New("user is empty")

func GetId(r *http.Request) (result int, err error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, IsEmpty
	}
	return strconv.Atoi(idStr)
}

func GetByName(r *http.Request, name string) (result string, err error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return "", IsEmpty
	}
	return value, nil
}

func GetUser(r *http.Request) (user *domain.User, err error) {
	err = json.NewDecoder(r.Body).Decode(&user)
	return user, err
}
