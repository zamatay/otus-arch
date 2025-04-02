package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"githib.com/zamatay/otus/arch/lesson-1/internal/app"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer done()

	config, err := app.NewConfig()
	if err != nil {
		log.Fatal("Ошибка при инициализации конфига приложения", err)
	}
	slog.Info("Загрузили конфиг", "config", config)

	repo, cache, producer, service, err := app.NewInfra(ctx, config)
	if err != nil {
		log.Fatal("Ошибка при инициализации приложения", err)
	}

	app.RegisterApi(repo, service, cache, producer, config.App.Secret)

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
