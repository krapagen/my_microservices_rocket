package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/krapagen/my_microservices_rocket/payment/internal/app"
	"github.com/krapagen/my_microservices_rocket/payment/internal/config"
)

func main() {
	// Загружаем переменные окружения из payment.env (если файл существует)
	err := godotenv.Load("../payment.env")
	if err != nil {
		slog.Warn("ошибка загрузки переменных из окружения .env", "error", err)
	}

	configPath := config.ResolveConfigPath()

	config.MustLoad(configPath)

	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
