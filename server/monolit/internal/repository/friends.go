package repository

import (
	"context"
	"fmt"
	"log/slog"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (r *Repo) SetFriends(ctx context.Context, fromUserId, toUserId int) (bool, error) {
	cmd, err := r.GetWriteConnection().Exec(ctx, `insert into friends(to_user_id, from_user_id) values($1, $2) ON CONFLICT(to_user_id, from_user_id) DO NOTHING`, toUserId, fromUserId)
	if err != nil {
		slog.Error("Ошибка при добавлении данных друзей", "error", err)
		return false, fmt.Errorf("Internal error")
	}
	return cmd.RowsAffected() != 0, nil
}

func (r *Repo) DeleteFriends(ctx context.Context, fromUserId, toUserId int) error {
	userFrom := domain.GetUserFromContext(ctx)
	if userFrom == nil {
		return fmt.Errorf("User not found")
	}
	_, err := r.GetWriteConnection().Exec(ctx, `delete from friends where to_user_id=$1 and from_user_id=$2`, toUserId, fromUserId)
	if err != nil {
		slog.Error("Ошибка при удалении данных друзей", "error", err)
		return fmt.Errorf("Internal error")
	}
	return nil
}
