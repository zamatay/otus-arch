package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// Структура для хранения данных пользователя из токена
type UserClaims struct {
	Username string `json:"username,omitempty"`
	Id       int    `json:"id,omitempty"`
	Exp      int64  `json:"exp,omitempty"`
	jwt.RegisteredClaims
}

func GetUserFromContext(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value("auth").(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}
