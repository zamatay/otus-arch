package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
)

const postFields = "id, user_id, text"

var internalError = fmt.Errorf("internal error")

func (r *Repo) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	row := r.GetWriteConnection().QueryRow(ctx,
		`insert into posts(user_id, text) 
			values ($1, $2)
			RETURNING `+postFields+`;`, post.UserID, post.Text)

	postObject, err := r.scanRow(row)

	return postObject, err
}

func (r *Repo) DeletePost(ctx context.Context, id string, userId int) (bool, error) {
	cmd, err := r.GetWriteConnection().Exec(ctx, `delete from posts where id=$1 and user_id=$2`, id, userId)
	if err != nil {
		slog.Error("Ошибка DeletePost", "error", err)
		return false, internalError
	}

	return cmd.RowsAffected() > 0, nil
}

func (r *Repo) UpdatePost(ctx context.Context, post *domain.Post) (bool, error) {
	cmd, err := r.GetWriteConnection().Exec(ctx, `update posts set text = $2 where id=$1 and user_id=$3`, post.ID, post.Text)
	if err != nil {
		slog.Error("Ошибка UpdatePost", "error", err)
		return false, internalError
	}
	return cmd.RowsAffected() > 0, nil
}

func (r *Repo) scanRow(row pgx.Row) (*domain.Post, error) {
	post := domain.Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.Text)
	if err != nil {
		slog.Error("Ошибка GetPost", "error", err)
		return nil, internalError
	}
	return &post, nil
}

func (r *Repo) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	row := r.GetConnection().QueryRow(ctx, `select `+postFields+` from posts where id=$1`, id)
	return r.scanRow(row)
}

func (r *Repo) FeedPost(ctx context.Context, offset int, limit int, userId int) ([]*domain.Post, error) {
	//userFrom := domain.GetUserFromContext(ctx)
	rows, err := r.GetConnection().Query(ctx, `select id, user_id, text from posts where user_id=$1 order by created_at desc limit $2 offset $3`, userId, limit, offset)
	if err != nil {
		slog.Error("Ошибка FeedPost", "error", err)
		return nil, internalError
	}
	posts := make([]*domain.Post, 0, 100)
	for rows.Next() {
		post := domain.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Text)
		if err != nil {
			slog.Error("Ошибка FeedPost", "error", err)
			return nil, internalError
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *Repo) Read(ctx context.Context, postId int, userId int) (int64, error) {
	_, err := r.GetWriteConnection().Exec(ctx, `insert into user_reade(user_id, post_id) values ($1, $2)
 		on conflict(user_id,post_id)
 		do update set update_ad = now()`, userId, postId)
	if err != nil {
		slog.Error("Ошибка Read", "error", err)
		return 0, internalError
	}
	row := r.GetConnection().QueryRow(ctx, `select count(*) from user_reade where post_id=$1`, postId)
	count := int64(0)
	if err = row.Scan(&count); err != nil {
		slog.Error("Ошибка Scan", "error", err)
		return 0, internalError
	}
	return count, nil
}
