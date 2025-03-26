package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

const _timeout time.Duration = 200 * time.Second

var IsEmpty = errors.New("Param is empty")

type OkResult struct{ Ok bool }
type CustomResult struct {
	OkResult
	Message string
}

func SetError(writer http.ResponseWriter, errorString string, code int) {
	http.Error(writer, errorString, code)
}

func Ok() OkResult {
	return OkResult{Ok: true}
}
func OkFalse(isOk bool) OkResult {
	return OkResult{Ok: isOk}
}
func SetOk(writer http.ResponseWriter, object any) {
	err := json.NewEncoder(writer).Encode(object)
	if err != nil {
		SetError(writer, err.Error(), 500)
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(http.StatusOK)
}

func GetContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, _timeout)

}

func GetByName(r *http.Request, name string) (result int, err error) {
	idStr := r.URL.Query().Get(name)
	if idStr == "" {
		return 0, IsEmpty
	}
	return strconv.Atoi(idStr)
}
