package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/app"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/config"
)

func main() {
	// Загружаем переменные окружения из ufo.env (если файл существует)
	err := godotenv.Load("../inventory.env")
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
