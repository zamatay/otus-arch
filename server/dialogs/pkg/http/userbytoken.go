package http

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"dialogs/internal/model"
)

var UserByToken func(tokenString string) (*model.UserClaims, error)

func SetUserByToken(secretKey []byte) func(tokenString string) (*model.UserClaims, error) {
	return func(tokenString string) (*model.UserClaims, error) {
		token, err := jwt.ParseWithClaims(tokenString, &model.UserClaims{}, func(token *jwt.Token) (any, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return nil, errors.New("Invalid token")
		}

		user, ok := token.Claims.(*model.UserClaims)
		if !ok {
			return nil, errors.New("Invalid token")
		}
		return user, nil
	}
}
