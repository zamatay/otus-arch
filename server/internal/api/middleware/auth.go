package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
		// Проверка токена
		token, err := jwt.Parse(a[1], func(token *jwt.Token) (any, error) {
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

		// Если токен действителен, передаем управление следующему обработчику
		next(w, r)
	}
}
