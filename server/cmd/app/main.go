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
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	config, err := app.NewConfig()
	if err != nil {
		log.Fatal("Ошибка при инициализации конфига приложения", err)
	}
	repo, err := repository.NewRepo(ctx, config.DB["read"], config.DB["write"])
	if err != nil {
		log.Fatal("Ошибка при инициализации репозитория", err)
	}

	service, err := api.New(&config.Http, config.App.Secret)
	if err != nil {
		return
	}
	user.NewUser(repo, service)
	auth.NewAuth(repo, service, config.App.Secret)
	friend.NewFriend(repo, service)
	post.NewPost(repo, service)
	if err := service.Start(); err != nil {
		log.Fatal("Ошибка при запуске http", err)
	}

	slog.Info("Приложение запустилось", slog.Uint64("port", uint64(config.Http.Port)))
	slog.Info("Поступил сигнал на завершение", slog.Any("done", <-ctx.Done()))

	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()

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
