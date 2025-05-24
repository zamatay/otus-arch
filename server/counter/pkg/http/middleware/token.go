package middleware

import (
	"context"
	"net/http"
	"strings"

	httpInternal "counter/pkg/http"
)

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

		user, err := httpInternal.UserByToken(a[1])

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return

		}

		ctx := context.WithValue(r.Context(), "auth", user)

		// Если токен действителен, передаем управление следующему обработчику
		next(w, r.WithContext(ctx))
	}
}
