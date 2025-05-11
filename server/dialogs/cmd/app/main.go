package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"dialogs/internal/app"
	"dialogs/internal/config"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer done()

	newConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal("Ошибка при инициализации конфига приложения", err)
	}

	repo, service, err := app.NewInfra(ctx, newConfig)
	if err != nil {
		log.Fatal("Ошибка при инициализации приложения", err)
	}

	app.RegisterApi(ctx, repo, service)

	if err := service.Start(); err != nil {
		log.Fatal("Ошибка при запуске http", err)
	}

	slog.Info("Приложение запустилось", slog.Uint64("port", uint64(newConfig.Http.Port)))
	slog.Info("Поступил сигнал на завершение", slog.Any("done", <-ctx.Done()))

}
