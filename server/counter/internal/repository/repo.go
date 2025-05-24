package repository

import (
	"context"
	"log"

	"counter/internal/repository/redis"
)

type Repo struct {
	cache *redis.Cache
}

func NewRepo(ctx context.Context, config *redis.Config) (repo *Repo) {
	repo = new(Repo)
	var err error
	repo.cache, err = redis.NewCache(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func (repo *Repo) Increase(ctx context.Context, postId int32) error {
	return repo.cache.Increase(ctx, postId)
}

func (repo *Repo) Decrease(ctx context.Context, postId int32) error {
	return repo.cache.Decrease(ctx, postId)
}

func (repo *Repo) GetPostRead(ctx context.Context, postId int32) (int, error) {
	return repo.cache.GetPostRead(ctx, postId)
}
