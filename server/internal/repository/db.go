package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

func (r *Repo) Close(ctx context.Context) error {
	r.conn.Close()
	return nil
}

func NewRepo(ctx context.Context, cfg Config) (*Repo, error) {
	dsn := cfg.GetConnectionString()
	return NewRepoByStr(ctx, dsn)
}

func NewRepoByStr(ctx context.Context, dsn string) (*Repo, error) {
	connect, err := pgxpool.New(ctx, dsn)
	connect.Config().MaxConns = 1000
	if err != nil {
		return nil, err
	}
	if err := connect.Ping(ctx); err != nil {
		return nil, err
	}
	return &Repo{
		conn: connect,
	}, nil
}
