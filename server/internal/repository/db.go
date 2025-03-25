package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	balancer  *RandomBalancer
	writeConn *pgxpool.Pool
}

func (r *Repo) Close(ctx context.Context) error {
	for _, pool := range r.balancer.GetAllReplica() {
		pool.Close()
	}
	return nil
}

func NewRepo(ctx context.Context, cfg []*Config, writeConfig []*Config) (*Repo, error) {
	repo := new(Repo)
	repo.balancer = NewRandomBalancer()
	for _, config := range cfg {
		dsn := config.GetConnectionString()
		connection, err := NewConnection(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("Ошибка подключение к базе данных: %w", err)
		}
		repo.balancer.AddReplica(connection)
	}

	if len(writeConfig) == 0 {
		return nil, fmt.Errorf("Не указаны конфигурации для записи")
	}

	for _, config := range writeConfig {
		dsn := config.GetConnectionString()
		connection, err := NewConnection(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("Ошибка подключение к базе данных: %w", err)
		}
		repo.writeConn = connection
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

func (r *Repo) GetConnection() *pgxpool.Pool {
	return r.balancer.GetReplica()
}

func (r *Repo) GetWriteConnection() *pgxpool.Pool {
	return r.writeConn
}
