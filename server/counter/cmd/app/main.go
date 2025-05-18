package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"counter/internal/app"
	"counter/internal/config"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer done()

	newConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal("Ошибка при инициализации конфига приложения", err)
	}

	repo := app.NewInfra(ctx, newConfig)

	service, grpc, err := app.NewService(newConfig, repo)
	if err != nil {
		log.Fatal("Ошибка при инициализации сервисов", err)
	}

	err = grpc.Register(ctx, service)
	if err != nil {
		log.Fatal("Ошибка при инициализации сервисов", err)
	}

	if err := service.Start(); err != nil {
		log.Fatal("Ошибка при запуске http", err)
	}

	grpc.Start()
	//
	slog.Info("Приложение запустилось", slog.Uint64("port", uint64(newConfig.Http.Port)))
	slog.Info("Поступил сигнал на завершение", slog.Any("done", <-ctx.Done()))

	deadlineContext, doneDeadline := context.WithTimeout(ctx, 5*time.Second)
	err = service.Stop(deadlineContext)
	if err != nil {
		slog.Error("Ошибка при завершении программы", "error", err)
	}
	doneDeadline()
}
