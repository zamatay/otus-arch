package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	RDB *redis.Client
}

func NewCache(ctx context.Context, config *Config) (*Cache, error) {
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

func (receiver *Cache) Decrease(ctx context.Context, postId int32) error {
	result, err := receiver.RDB.Decr(ctx, fmt.Sprintf("post:couner:%d", postId)).Result()
	if err != nil {
		return err
	}
	if result < 0 {
		_, err := receiver.RDB.Set(ctx, fmt.Sprintf("post:couner:%d", postId), 0, 0).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (receiver *Cache) Increase(ctx context.Context, postId int32) error {
	_, err := receiver.RDB.Incr(ctx, fmt.Sprintf("post:couner:%d", postId)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Cache) GetPostRead(ctx context.Context, postId int32) (int, error) {
	countStr, err := receiver.RDB.Get(ctx, fmt.Sprintf("post:couner:%d", postId)).Result()
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}
	return count, nil
}
