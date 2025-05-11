package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zamatay/otus/arch/lesson-1/internal/kafka"
)

type Repo struct {
	balancer  *RandomBalancer
	writeConn *pgxpool.Pool
	shardConn *pgxpool.Pool
	Producer  *kafka.Producer
}

func (r *Repo) Close(ctx context.Context) error {
	for _, pool := range r.balancer.GetAllReplica() {
		pool.Close()
	}
	r.writeConn.Close()
	return nil
}

func NewRepo(ctx context.Context, cfg []*Config, writeConfig []*Config, shardConfig []*Config) (*Repo, error) {
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

	if len(shardConfig) == 0 {
		return nil, fmt.Errorf("Не указаны конфигурации для шарда")
	}

	for _, config := range shardConfig {
		dsn := config.GetConnectionString()
		connection, err := NewConnection(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("Ошибка подключение к базе данных: %w", err)
		}
		repo.shardConn = connection
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
func (r *Repo) GetShardConnection() *pgxpool.Pool {
	return r.shardConn
}
