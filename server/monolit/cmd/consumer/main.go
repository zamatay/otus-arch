package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/IBM/sarama"

	"github.com/zamatay/otus/arch/lesson-1/internal/config"
	"github.com/zamatay/otus/arch/lesson-1/internal/domain"
	"github.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

const (
	cacheCount = 1000
	timeout    = 500 * time.Millisecond
)

type Cache struct {
	cache *redis.Cache
}

func NewCache(cfg *redis.Config) (*Cache, error) {
	cache, err := redis.NewCache(context.Background(), *cfg)
	if err != nil {
		return nil, err
	}
	return &Cache{cache: cache}, nil
}

func main() {
	configKafka := sarama.NewConfig()
	configKafka.Consumer.Return.Errors = true

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cache, err := NewCache(&cfg.Cache)

	chanel := kafka.GetChannel()

	newConsumer, err := kafka.NewConsumer(&cfg.Kafka, chanel)
	if err != nil {
		log.Fatal(err)
	}
	defer newConsumer.Close()

	go func() {
		for msg := range chanel {
			for _, header := range msg.Headers {
				if string(header.Key) != "message-type" {
					continue
				}

				switch string(msg.Value) {
				case "create":
					cache.addToCache(msg)
				case "delete":
					cache.deleteFromCache(msg)
				}

			}
		}

	}()

	newConsumer.Process()
}

func (r *Cache) deleteFromCache(msg *sarama.ConsumerMessage) {
	ctx, done := context.WithTimeout(context.Background(), timeout)
	defer done()

	deletePost, err := getPost(msg)
	if err != nil {
		slog.Error("Ошибка при получении поста из сообщения", "error", err, "message", msg.Value)
		return
	}

	needDel := false
	key := r.cache.GetFeedKey(deletePost.UserID)
	result, err := r.cache.RDB.LRange(ctx, key, 0, cacheCount).Result()
	if err != nil {
		slog.Error("Ошибка DeletePost", "error", err)
		return
	}
	posts := make([]*domain.Post, 0, len(result))
	for _, item := range result {
		currentPost := domain.Post{}
		err := json.Unmarshal([]byte(item), &currentPost)
		if err != nil {
			return
		}
		if currentPost.ID == deletePost.ID {
			needDel = true
			continue
		}
		posts = append(posts, &currentPost)
	}
	if needDel {
		r.cache.RDB.Del(ctx, key)
		err := r.cache.InsertPosts(ctx, posts)
		if err != nil {
			slog.Error("Ошибка при добавлении Posts", "error", err)
		}
	}
}

func (r *Cache) addToCache(message *sarama.ConsumerMessage) {
	ctx, done := context.WithTimeout(context.Background(), timeout)
	defer done()

	p, err := getPost(message)
	if err != nil {
		slog.Error("Ошибка при получении поста из сообщения", "error", err, "message", message.Value)
		return
	}
	userLenPost := r.cache.GetLen(ctx, r.cache.GetFeedKey(p.UserID))
	if userLenPost >= cacheCount {
		countDel := userLenPost - cacheCount - 1
		if err := r.cache.DelLeftFeedCache(ctx, p.UserID, int(countDel)); err != nil {
			slog.Error("Ошибка при удалении записи из кэша", "error", err)
		}
	}

	err = r.cache.SetFeedCache(ctx, p.UserID, []*domain.Post{p})
	if err != nil {
		slog.Error("Ошибка при добавлении записи в кэш", "error", err)
	}

}

func getPost(message *sarama.ConsumerMessage) (*domain.Post, error) {
	msg := string(message.Value)
	p := domain.Post{}
	if err := json.Unmarshal(message.Value, &p); err != nil {
		slog.Error("Ошибка при сериализации сообщения", "error", err, "message", msg)
		return nil, err
	}
	return &p, nil
}
