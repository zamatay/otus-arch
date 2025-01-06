package domain

type User struct {
	ID        int      `json:"id"`
	Login     string   `json:"login"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Birthday  string   `json:"birthday"`
	GenderID  int      `json:"gender_id"`
	Interests []string `json:"interests"`
	City      string   `json:"city"`
	Enabled   bool     `json:"enabled"`
}

type UserCredentials struct {
	UserID       int    `json:"user_id"`
	PasswordHash string `json:"password_hash"`
}
