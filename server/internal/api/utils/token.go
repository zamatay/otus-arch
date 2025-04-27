package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

var UserByToken func(tokenString string) (*domain.UserClaims, error)

func SetUserByToken(secretKey []byte) func(tokenString string) (*domain.UserClaims, error) {
	return func(tokenString string) (*domain.UserClaims, error) {
		token, err := jwt.ParseWithClaims(tokenString, &domain.UserClaims{}, func(token *jwt.Token) (any, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return nil, errors.New("Invalid token")
		}

		user, ok := token.Claims.(*domain.UserClaims)
		if !ok {
			return nil, errors.New("Invalid token")
		}
		return user, nil
	}
}
