package tests

import (
	"testing"
	"time"

	"otus-arch/server/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    domain.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: domain.User{
				ID:        1,
				Login:     "testuser",
				FirstName: "Test",
				LastName:  "User",
				Birthday:  time.Now().Add(-20 * 365 * 24 * time.Hour),
				GenderID:  1,
				Interests: []string{"reading", "coding"},
				City:      "Moscow",
				Enabled:   true,
			},
			wantErr: false,
		},
		{
			name: "empty login",
			user: domain.User{
				ID:        1,
				Login:     "",
				FirstName: "Test",
				LastName:  "User",
				Birthday:  time.Now().Add(-20 * 365 * 24 * time.Hour),
				GenderID:  1,
				Interests: []string{"reading", "coding"},
				City:      "Moscow",
				Enabled:   true,
			},
			wantErr: true,
		},
		{
			name: "future birthday",
			user: domain.User{
				ID:        1,
				Login:     "testuser",
				FirstName: "Test",
				LastName:  "User",
				Birthday:  time.Now().Add(24 * time.Hour),
				GenderID:  1,
				Interests: []string{"reading", "coding"},
				City:      "Moscow",
				Enabled:   true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_IsAdult(t *testing.T) {
	tests := []struct {
		name     string
		birthday time.Time
		want     bool
	}{
		{
			name:     "adult user",
			birthday: time.Now().Add(-20 * 365 * 24 * time.Hour),
			want:     true,
		},
		{
			name:     "minor user",
			birthday: time.Now().Add(-15 * 365 * 24 * time.Hour),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := domain.User{
				Birthday: tt.birthday,
			}
			got := user.IsAdult()
			assert.Equal(t, tt.want, got)
		})
	}
}
