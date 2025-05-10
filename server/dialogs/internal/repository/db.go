package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

func NewRepo(ctx context.Context, config *Config) (repo *Repo, err error) {
	repo = new(Repo)
	dsn := config.GetConnectionString()
	if repo.conn, err = NewConnection(ctx, dsn); err != nil {
		return nil, fmt.Errorf("Ошибка подключение к базе данных: %w", err)
	}
	return repo, nil
}

func NewConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	connect, err := pgxpool.New(ctx, dsn)
	connect.Config().MaxConns = 1000
	if err != nil {
		return nil, err
	}
	if err := connect.Ping(ctx); err != nil {
		return nil, err
	}
	return connect, nil
}
