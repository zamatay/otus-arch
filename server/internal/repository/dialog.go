package repository

import (
	"context"
	"fmt"
	"log/slog"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
)

func (r *Repo) SendDialog(ctx context.Context, fromUserId int, toUserId int, text string) (bool, error) {
	cmd, err := r.GetShardConnection().Exec(ctx, `insert into dialogs(from_user_id, to_user_id, text) values($1, $2, $3)`, fromUserId, toUserId, text)
	if err != nil {
		slog.Error("Ошибка при добавлении диалога", "error", err)
		return false, fmt.Errorf("Internal error")
	}
	return cmd.RowsAffected() != 0, nil
}
func (r *Repo) ListDialog(ctx context.Context, fromUserId int, toUserId int) ([]*domain.Dialog, error) {
	rows, err := r.GetShardConnection().Query(ctx, `Select from_user_id, to_user_id, text, created_at, updated_at FROM dialogs WHERE from_user_id = $1 AND to_user_id = $2`, fromUserId, toUserId)
	if err != nil {
		slog.Error("Ошибка при получении диалога", "error", err)
		return nil, fmt.Errorf("Internal error")
	}
	dialogs := make([]*domain.Dialog, 0)
	for rows.Next() {
		var dialog domain.Dialog
		err := rows.Scan(&dialog.FromUserID, &dialog.ToUserID, &dialog.Text, &dialog.CreatedAt, &dialog.UpdatedAt)
		if err != nil {
			slog.Error("Ошибка при сканирование диалога", "error", err)
			return nil, err
		}
		dialogs = append(dialogs, &dialog)
	}

	return dialogs, nil
}
