package auth

import (
	"log/slog"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

const _expired = time.Hour * 1

func CreateToken(u domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": u.Login,
		"id":       u.ID,
		"exp":      time.Now().Add(_expired).Unix(),
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func HashPassword(password string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Ошибка получения hash пароля", "error", err)
		return ""
	}
	return string(hashPassword)

}

func ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Error("Ошибка получения hash пароля", "error", err)
		return false
	}
	return true

}
