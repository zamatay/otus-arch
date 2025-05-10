package repository

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/zamatay/otus/arch/lesson-1/internal/api/auth"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (r *Repo) Login(ctx context.Context, name string, password string) (string, *domain.User, error) {
	user := r.GetUserIdByLogin(ctx, name)
	if user == nil {
		return "", nil, errors.New("Пользователь не найден")
	}

	row := r.GetConnection().QueryRow(ctx, "Select password_hash from user_credentials where user_id = $1", user.ID)
	var passwordHash string
	if err := row.Scan(&passwordHash); err != nil {
		return "", nil, err
	}

	if !auth.ComparePassword(passwordHash, password) {
		return "", nil, errors.New("Некорректный пароль")
	}

	token, err := auth.CreateToken(*user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (r *Repo) Register(ctx context.Context, user domain.RegisterUser) error {
	if !PasswordIsValid(user.Password) {
		return errors.New("Пароль не прошел валидацию")
	}

	tx, err := r.GetWriteConnection().BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	//Добавляем пользователя в таблицу user
	row := tx.QueryRow(ctx, `insert into users(login, first_name, last_name, birthday, gender_id, city, enabled, interests)
		values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		user.Login, user.FirstName, user.LastName, user.Birthday, user.GenderID, user.City, true, user.Interests)
	var id int
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(context.Background()); err != nil {
			slog.Error("Ошибка при откате транзакции", "error", err)
		}
		slog.Error("Ошибка при добавление пользователя в таблицу user", "error", err)
		return err
	}
	//Добавляем пользователя в таблицу user_credentials
	hash := auth.HashPassword(user.Password)
	if _, err = tx.Exec(ctx, `insert into user_credentials(user_id, password_hash)
		values ($1,$2)`, id, hash); err != nil {
		if err := tx.Rollback(context.Background()); err != nil {
			slog.Error("Ошибка при откате транзакции", "error", err)
		}
		slog.Error("Ошибка при добавление пользователя в таблицу user_credentials", "error", err)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при commit транзакции", "error", err)
		return err
	}
	slog.Info("Пользователь создан", "user", user)
	return nil
}

func PasswordIsValid(password string) bool {
	return len(strings.Trim(password, " ")) >= 6
}
