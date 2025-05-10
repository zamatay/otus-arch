package tests

import (
	"context"
	"testing"
	"time"

	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*repository.Repo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	conn, err := pgx.NewConn(pgx.ConnConfig{
		Config: db,
	})
	require.NoError(t, err)

	repo := &repository.Repo{
		conn: conn,
	}

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestRepo_GetUser(t *testing.T) {
	repo, mock, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	expectedUser := domain.User{
		ID:        1,
		Login:     "testuser",
		FirstName: "Test",
		LastName:  "User",
		Birthday:  time.Now(),
		GenderID:  1,
		City:      "Moscow",
		Enabled:   true,
		Interests: []string{"reading", "coding"},
	}

	rows := sqlmock.NewRows([]string{"id", "login", "first_name", "last_name", "birthday", "gender_id", "city", "enabled", "interests"}).
		AddRow(expectedUser.ID, expectedUser.Login, expectedUser.FirstName, expectedUser.LastName, expectedUser.Birthday, expectedUser.GenderID, expectedUser.City, expectedUser.Enabled, expectedUser.Interests)

	mock.ExpectQuery("select (.+) from users where id=\\$1").
		WithArgs(1).
		WillReturnRows(rows)

	user := repo.GetUser(ctx, 1)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Login, user.Login)
	assert.Equal(t, expectedUser.FirstName, user.FirstName)
	assert.Equal(t, expectedUser.LastName, user.LastName)
	assert.Equal(t, expectedUser.City, user.City)
	assert.Equal(t, expectedUser.Enabled, user.Enabled)
	assert.Equal(t, expectedUser.Interests, user.Interests)
}

func TestRepo_AddUser(t *testing.T) {
	repo, mock, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := domain.User{
		Login:     "newuser",
		FirstName: "New",
		LastName:  "User",
		Birthday:  time.Now(),
		GenderID:  1,
		City:      "Moscow",
		Enabled:   true,
	}

	mock.ExpectQuery("insert into users").
		WithArgs(user.Login, user.FirstName, user.LastName, user.Birthday, user.GenderID, user.City, user.Enabled).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.AddUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestRepo_UpdateUser(t *testing.T) {
	repo, mock, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	user := domain.User{
		ID:        1,
		Login:     "updateduser",
		FirstName: "Updated",
		LastName:  "User",
		Birthday:  time.Now(),
		GenderID:  1,
		City:      "Moscow",
		Enabled:   true,
		Interests: []string{"reading"},
	}

	mock.ExpectExec("update users").
		WithArgs(user.ID, user.FirstName, user.LastName, user.Birthday, user.GenderID, user.City, user.Enabled, user.Interests).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateUser(ctx, user)
	assert.NoError(t, err)
}

func TestRepo_Remove(t *testing.T) {
	repo, mock, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectExec("delete from users where id=\\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Remove(ctx, 1)
	assert.NoError(t, err)
}

func TestRepo_SearchUser(t *testing.T) {
	repo, mock, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	expectedUsers := []domain.User{
		{
			ID:        1,
			Login:     "testuser1",
			FirstName: "Test",
			LastName:  "User1",
			Birthday:  time.Now(),
			GenderID:  1,
			City:      "Moscow",
			Enabled:   true,
		},
		{
			ID:        2,
			Login:     "testuser2",
			FirstName: "Test",
			LastName:  "User2",
			Birthday:  time.Now(),
			GenderID:  1,
			City:      "Moscow",
			Enabled:   true,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "login", "first_name", "last_name", "birthday", "gender_id", "city", "enabled", "interests"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.Login, user.FirstName, user.LastName, user.Birthday, user.GenderID, user.City, user.Enabled, user.Interests)
	}

	mock.ExpectQuery("select (.+) from users where left\\(first_name, 3\\) = left\\(\\$1, 3\\) and left\\(last_name, 3\\) = left\\(\\$2, 3\\) limit 100").
		WithArgs("Test", "User").
		WillReturnRows(rows)

	users, err := repo.SearchUser(ctx, "Test", "User")
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].ID, users[0].ID)
	assert.Equal(t, expectedUsers[1].ID, users[1].ID)
}
