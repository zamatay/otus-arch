package middleware

import (
	"net/http"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

func CorsMiddleware(next HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Разрешить все источники
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Разрешенные методы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Разрешенные заголовки

		// Обработка preflight-запроса
		if r.Method == http.MethodOptions {
			return
		}

		next(w, r.WithContext(r.Context()))
	}
}
