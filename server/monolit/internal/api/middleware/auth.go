package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zamatay/otus/arch/lesson-1/internal/api/utils"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

// Middleware для проверки токена
func TokenMiddleware(next HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Строка авторизации пуста", http.StatusUnauthorized)
			return
		}

		a := strings.Split(tokenString, " ")

		if len(a) != 2 {
			http.Error(w, "Invalid authorization", http.StatusUnauthorized)
			return

		}

		user, err := utils.UserByToken(a[1])

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return

		}

		ctx := context.WithValue(r.Context(), "auth", user)

		// Если токен действителен, передаем управление следующему обработчику
		next(w, r.WithContext(ctx))
	}
}

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
