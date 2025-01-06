package repository

import (
	"context"
	"log/slog"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (r Repo) GetUsers(ctx context.Context) []domain.User {
	rows, err := r.conn.Query(ctx, "select id, login, first_name, last_name, birthday, gender_id, city, enabled from users limit 100")
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]domain.User, 0, 100)
	for rows.Next() {
		u := domain.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.Birthday, &u.GenderID, &u.City, &u.Enabled)
		if err != nil {
			slog.Error("Ошибка при сканировании данных в структуру")
			continue
		}
		result = append(result, u)
	}
	return result
	//return []domain.User{{ID: 1, Login: "zamatay", FirstName: "Александр", LastName: "Замураев", Birthday: "12-03-1978", GenderID: 1, City: "Краснодар", Enabled: true}}
}

func (r Repo) GetUser(ctx context.Context, id int) *domain.User {
	row := r.conn.QueryRow(ctx, "select id, login, first_name, last_name, birthday, gender_id, city, enabled from users where id=$1", id)

	u := domain.User{}
	err := row.Scan(&u.ID, &u.Login, &u.FirstName, &u.LastName, &u.Birthday, &u.GenderID, &u.City, &u.Enabled)
	if err != nil {
		slog.Error("Ошибка при сканировании данных в структуру")
		return nil
	}

	return &u
}

func (r Repo) AddUser(ctx context.Context, user domain.User) (int, error) {
	r.conn.Exec(ctx, "insert into users(id, login, first_name, last_name, birthday, gender_id, city, enabled) "+
		"values()")
}

func (r Repo) UpdateUser(ctx context.Context, i int, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Remove(ctx context.Context, i int) error {
	//TODO implement me
	panic("implement me")
}
