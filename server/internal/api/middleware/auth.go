package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

// Middleware для проверки токена
func TokenMiddleware(next HandleFunc, secretKey []byte) HandleFunc {
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

		token, err := jwt.ParseWithClaims(a[1], &domain.UserClaims{}, func(token *jwt.Token) (any, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, ok := token.Claims.(*domain.UserClaims)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "auth", user)

		// Если токен действителен, передаем управление следующему обработчику
		next(w, r.WithContext(ctx))
	}
}
