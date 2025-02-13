package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func (r *Repo) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func NewRepo(ctx context.Context, cfg Config) (*Repo, error) {
	dsn := cfg.GetConnectionString()
	return NewRepoByStr(ctx, dsn)
}

func NewRepoByStr(ctx context.Context, dsn string) (*Repo, error) {
	connect, err := pgx.Connect(ctx, dsn)
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
