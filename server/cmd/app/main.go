package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/app"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	config, err := app.NewConfig()
	if err != nil {
		log.Fatal("Ошибка при инициализации конфига приложения", err)
	}
	repo, err := repository.NewRepo(ctx, config.DB)
	if err != nil {
		log.Fatal("Ошибка при инициализации репозитория", err)
	}
	_ = repo

	api.NewUser(repo, config.Http)

	slog.Info("Приложение запустилось", slog.Uint64("порт", uint64(config.Http.Port)))
	slog.Info("Приложение завершилось", slog.Any("done", <-ctx.Done()))

	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()

	if err := repo.Close(ctx); err != nil {
		log.Fatal("Ошибка при закрытии БД", err)
	}
	slog.Info("База данных закрылась")
}
