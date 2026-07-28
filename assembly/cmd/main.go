package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/app"
	"github.com/krapagen/my_microservices_rocket/assembly/internal/config"
)

func main() {
	// Загружаем переменные окружения из assembly.env (если файл существует)
	err := godotenv.Load("../assembly.env")
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
