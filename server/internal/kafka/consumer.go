package kafka

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"

	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

const (
	cacheCount = 1000
	timeout    = 500 * time.Millisecond
)

type Consumer struct {
	config         *sarama.Config
	instance       sarama.Consumer
	partitionPosts sarama.PartitionConsumer
	cache          *redis.Cache
}

func NewConsumer(cfg *Config, cache *redis.Cache) (*Consumer, error) {
	kafka := new(Consumer)

	kafka.config = sarama.NewConfig()
	kafka.config.Consumer.Return.Errors = true
	kafka.cache = cache

	var err error
	if kafka.instance, err = sarama.NewConsumer([]string{cfg.GetAddress()}, kafka.config); err != nil {
		return nil, err
	}

	kafka.partitionPosts, err = kafka.instance.ConsumePartition("posts", 0, sarama.OffsetNewest)

	return kafka, nil
}

func (k *Consumer) Close() error {
	err := k.partitionPosts.Close()
	err = k.instance.Close()
	return err
}

func (k *Consumer) Process() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-k.partitionPosts.Messages():
			log.Printf("Received message: %s\n", string(msg.Value))

			for _, header := range msg.Headers {
				if string(header.Key) != "message-type" {
					continue
				}
				switch string(header.Value) {
				case "create":
					k.addToCache(msg)
				case "delete":
					k.deleteFromCache(msg)
				}
			}

		case err := <-k.partitionPosts.Errors():
			log.Println("Consumer error:", err)

		case <-signals:
			slog.Info("Приложение завершило работу")
			return
		}
	}

}

func (k *Consumer) deleteFromCache(msg *sarama.ConsumerMessage) {
	ctx, done := context.WithTimeout(context.Background(), timeout)
	defer done()

	deletePost, err := getPost(msg)
	if err != nil {
		slog.Error("Ошибка при получении поста из сообщения", "error", err, "message", msg.Value)
		return
	}

	needDel := false
	key := k.cache.GetFeedKey(deletePost.UserID)
	result, err := k.cache.RDB.LRange(ctx, key, 0, cacheCount).Result()
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
		k.cache.RDB.Del(ctx, key)
		err := k.cache.InsertPosts(ctx, posts)
		if err != nil {
			slog.Error("Ошибка при добавлении Posts", "error", err)
		}
	}
}

func (k *Consumer) addToCache(message *sarama.ConsumerMessage) {
	ctx, done := context.WithTimeout(context.Background(), timeout)
	defer done()

	p, err := getPost(message)
	if err != nil {
		slog.Error("Ошибка при получении поста из сообщения", "error", err, "message", message.Value)
		return
	}
	userLenPost := k.cache.GetLen(ctx, k.cache.GetFeedKey(p.UserID))
	if userLenPost >= cacheCount {
		countDel := userLenPost - cacheCount - 1
		if err := k.cache.DelLeftFeedCache(ctx, p.UserID, int(countDel)); err != nil {
			slog.Error("Ошибка при удалении записи из кэша", "error", err)
		}
	}

	err = k.cache.SetFeedCache(ctx, p.UserID, []*domain.Post{p})
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
