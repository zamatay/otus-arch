package repository

import (
	"context"
	"fmt"
	"log/slog"

	"dialogs/internal/model"
)

func (r *Repo) SendDialog(ctx context.Context, fromUserId int, toUserId int, text string) (*model.Dialog, error) {
	result := model.Dialog{}
	err := r.conn.QueryRow(ctx, `insert into dialogs(from_user_id, to_user_id, text) values($1, $2, $3) RETURNING from_user_id, to_user_id, text;`,
		fromUserId, toUserId, text).Scan(&result.FromUserID, &result.ToUserID, &result.Text)
	if err != nil {
		slog.Error("Ошибка при добавлении диалога", "error", err)
		return nil, fmt.Errorf("Internal error")
	}
	return &result, nil
}
func (r *Repo) ListDialog(ctx context.Context, fromUserId int, toUserId int) ([]*model.Dialog, error) {
	rows, err := r.conn.Query(ctx, `Select from_user_id, to_user_id, text, created_at, updated_at FROM dialogs WHERE from_user_id = $1 AND to_user_id = $2`, fromUserId, toUserId)
	if err != nil {
		slog.Error("Ошибка при получении диалога", "error", err)
		return nil, fmt.Errorf("Internal error")
	}
	dialogs := make([]*model.Dialog, 0)
	for rows.Next() {
		var dialog model.Dialog
		err := rows.Scan(&dialog.FromUserID, &dialog.ToUserID, &dialog.Text, &dialog.CreatedAt, &dialog.UpdatedAt)
		if err != nil {
			slog.Error("Ошибка при сканирование диалога", "error", err)
			return nil, err
		}
		dialogs = append(dialogs, &dialog)
	}

	return dialogs, nil
}
