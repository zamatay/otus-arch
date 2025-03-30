package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/auth"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/friend"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/post"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/user"
	"githib.com/zamatay/otus/arch/lesson-1/internal/app"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer done()

	config, err := app.NewConfig()
	slog.Info("Загрузили конфиг", "config", config)
	if err != nil {
		log.Fatal("Ошибка при инициализации конфига приложения", err)
	}

	repo, err := repository.NewRepo(ctx, config.DB["read"], config.DB["write"])
	if err != nil {
		log.Fatal("Ошибка при инициализации репозитория", err)
	}

	cache, err := redis.NewCache(ctx, config.Cache)

	producer, err := kafka.NewProducer(config.Kafka)
	service, err := api.New(&config.Http, config.App.Secret)
	if err != nil {
		return
	}
	user.NewUser(repo, service)
	auth.NewAuth(repo, service, config.App.Secret)
	friend.NewFriend(repo, service)
	post.NewPost(repo, cache, service, producer)
	if err := service.Start(); err != nil {
		log.Fatal("Ошибка при запуске http", err)
	}

	slog.Info("Приложение запустилось", slog.Uint64("port", uint64(config.Http.Port)))
	slog.Info("Поступил сигнал на завершение", slog.Any("done", <-ctx.Done()))

	ctx, doneTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer doneTimeout()

	slog.Info("Приложение начало закрываться")

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return repo.Close(ctx)
	})
	eg.Go(func() error {
		return service.Stop(ctx)
	})
	if err := eg.Wait(); err != nil {
		log.Fatal("Ошибка при закрытии приложения", err)
	}

	slog.Info("Приложение закрылось")
}
