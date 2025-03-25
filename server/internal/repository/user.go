package repository

import (
	"context"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

const selectUsers = "select id, login, first_name, last_name, birthday, gender_id, city, enabled, interests from users"

func (r *Repo) GetUsers(ctx context.Context) []domain.User {
	rows, err := r.GetConnection().Query(ctx, selectUsers+" limit 100")
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]domain.User, 0, 100)
	for rows.Next() {
		u := domain.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.Birthday, &u.GenderID, &u.City, &u.Enabled, &u.Interests)
		if err != nil {
			slog.Error("Ошибка при сканировании данных в структуру")
			continue
		}
		result = append(result, u)
	}
	return result
}

func (r *Repo) GetUser(ctx context.Context, id int) *domain.User {
	row := r.GetConnection().QueryRow(ctx, selectUsers+" where id=$1", id)
	return GetUserByRow(row)
}

func (r *Repo) GetUserIdByLogin(ctx context.Context, login string) *domain.User {
	row := r.GetConnection().QueryRow(ctx, selectUsers+` where login = $1`, login)
	return GetUserByRow(row)
}

func GetUserByRow(row pgx.Row) (u *domain.User) {
	u = new(domain.User)
	err := row.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.Birthday, &u.GenderID, &u.City, &u.Enabled, &u.Interests)
	if err != nil {
		slog.Error("Ошибка при сканировании данных в структуру")
		return nil
	}

	return u
}

func (r *Repo) AddUser(ctx context.Context, user domain.User) (int, error) {
	row := r.GetConnection().QueryRow(ctx, `insert into users(login, first_name, last_name, birthday, gender_id, city, enabled) 
		values($1,$2,$3,$4,$5,$6,$7) RETURNING id`, user.Login, user.FirstName, user.LastName, user.Birthday, user.GenderID, user.City, user.Enabled)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) UpdateUser(ctx context.Context, user domain.User) error {
	_, err := r.GetWriteConnection().Exec(ctx, `update users set first_name = $2, last_name = $3, birthday = $4, gender_id = $5, city = $6, enabled = $7, interests=$8 where id=$1`,
		user.ID, user.FirstName, user.LastName, user.Birthday, user.GenderID, user.City, user.Enabled, user.Interests)
	return err
}

func (r *Repo) Remove(ctx context.Context, id int) error {
	_, err := r.GetWriteConnection().Exec(ctx, "delete from users where id=$1", id)
	return err
}

func (r *Repo) SearchUser(ctx context.Context, firstName string, lastName string) ([]domain.User, error) {
	rows, err := r.GetConnection().Query(ctx, selectUsers+" where left(first_name, 3) = left($1, 3) and left(last_name, 3) = left($2, 3) limit 100", firstName, lastName)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, 0, 100)
	for rows.Next() {
		user := GetUserByRow(rows)
		if user != nil {
			users = append(users, *user)
		}
	}
	return users, nil
}

func addPercent(name string) string {
	if strings.HasSuffix(name, "%") {
		return name
	}
	return name + "%"
}
