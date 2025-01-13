package domain

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	GenderID  int       `json:"gender_id"`
	Interests []string  `json:"interests"`
	City      string    `json:"city"`
	Enabled   bool      `json:"enabled"`
}

type UserCredentials struct {
	UserID       int    `json:"user_id"`
	PasswordHash string `json:"password_hash"`
}

type AuthUser struct {
	User     string `json:"login"`
	Password string `json:"password"`
}

type AuthUserResult struct {
	Token string `json:"token"`
}

type RegisterUser struct {
	User
	Password string `json:"password"`
}
