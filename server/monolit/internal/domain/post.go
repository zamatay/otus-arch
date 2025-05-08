package domain

import (
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
