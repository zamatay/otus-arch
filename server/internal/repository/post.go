package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

const postFields = "id, user_id, text"

func (r *Repo) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	row := r.GetWriteConnection().QueryRow(ctx,
		`insert into posts(user_id, text) 
			values ($1, $2)
			RETURNING `+postFields+`;`, post.UserID, post.Text)
	return r.scanRow(row)
}

func (r *Repo) DeletePost(ctx context.Context, id int) (bool, error) {
	userFrom := domain.GetUserFromContext(ctx)
	cmd, err := r.GetWriteConnection().Exec(ctx, `delete from posts where id=$1 and user_id=$2`, id, userFrom.Id)
	if err != nil {
		slog.Error("Ошибка DeletePost", "error", err)
		return false, fmt.Errorf("Internal error")
	}
	return cmd.RowsAffected() > 0, nil
}

func (r *Repo) UpdatePost(ctx context.Context, post *domain.Post) (bool, error) {
	userFrom := domain.GetUserFromContext(ctx)
	cmd, err := r.GetWriteConnection().Exec(ctx, `update posts set text = $2 where id=$1 and user_id=$3`, post.ID, post.Text, userFrom.Id)
	if err != nil {
		slog.Error("Ошибка UpdatePost", "error", err)
		return false, fmt.Errorf("Internal error")
	}
	return cmd.RowsAffected() > 0, nil
}

func (r *Repo) scanRow(row pgx.Row) (*domain.Post, error) {
	post := domain.Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.Text)
	if err != nil {
		slog.Error("Ошибка GetPost", "error", err)
		return nil, fmt.Errorf("Internal error")
	}
	return &post, nil
}

func (r *Repo) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	row := r.GetConnection().QueryRow(ctx, `select `+postFields+` from posts where id=$1`, id)
	return r.scanRow(row)
}

func (r *Repo) FeedPost(ctx context.Context, offset, limit int) ([]*domain.Post, error) {
	userFrom := domain.GetUserFromContext(ctx)
	rows, err := r.GetConnection().Query(ctx, `select id, user_id, text from posts where user_id=$1 order by created_at desc limit $2 offset $3`, userFrom.Id, limit, offset)
	if err != nil {
		slog.Error("Ошибка FeedPost", "error", err)
		return nil, fmt.Errorf("Internal error")
	}
	posts := make([]*domain.Post, 0, 100)
	for rows.Next() {
		post := domain.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Text)
		if err != nil {
			slog.Error("Ошибка FeedPost", "error", err)
			return nil, fmt.Errorf("Internal error")
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
