package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	RDB *redis.Client
}

func NewCache(ctx context.Context, config Config) (*Cache, error) {
	cache := new(Cache)
	cache.RDB = redis.NewClient(&redis.Options{
		Addr:     config.GetAddress(),
		Username: config.User,
		Password: config.Password,
		DB:       config.DB,
	})

	err := cache.RDB.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return cache, nil
}
