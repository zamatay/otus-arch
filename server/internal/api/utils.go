package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const _timeout time.Duration = 200 * time.Second

type OkResult struct{ Ok bool }

func SetError(writer http.ResponseWriter, errorString string, code int) {
	http.Error(writer, errorString, code)
}

func Ok() OkResult {
	return OkResult{Ok: true}
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

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), _timeout)

}
