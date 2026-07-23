package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	inventoryApi "github.com/krapagen/my_microservices_rocket/inventory/internal/api/inventory/v1"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/config"
	inventoryRepository "github.com/krapagen/my_microservices_rocket/inventory/internal/repository/part"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/application/part"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

// diContainer — контейнер зависимостей (Composition Root) приложения
//
// Зачем это нужно:
// В простых приложениях зависимости создаются прямо в main.go: pool := pgxpool.New(...),
// repo := NewRepo(pool), svc := NewService(repo) и т.д. Это работает, пока зависимостей мало
// Когда сервис обрастает десятками компонентов, main.go превращается в «простыню» инициализации,
// а порядок создания начинает зависеть от неочевидных связей
//
// DI-контейнер решает эту проблему: каждый компонент «знает», от чего зависит, и создаёт
// свои зависимости по цепочке автоматически при первом обращении
//
// Как это работает:
// Каждый геттер (PGPool, inventoryRepo, inventoryService, inventoryV1Handler) следует паттерну
// «ленивая инициализация» (lazy initialization):
//  1. Проверяет, создан ли уже объект (nil-check)
//  2. Если нет — создаёт, запоминает в поле и возвращает
//  3. Если да — сразу возвращает ранее созданный экземпляр
//
// Это гарантирует, что каждый компонент создаётся ровно один раз, независимо от того,
// сколько раз к нему обращаются, и в правильном порядке
//
// Как добавить новую зависимость:
//  1. Добавьте поле с типом интерфейса в структуру
//  2. Напишите геттер с nil-check, который вызывает геттеры зависимостей
//  3. Используйте геттер там, где нужен компонент
//
// Почему интерфейсы (а не конкретные типы):
// Структуры слоёв (repository, service, api) — unexported, чтобы их нельзя было создать
// в обход конструктора New(). Контейнер хранит интерфейсы, которые определены в потребителях
// (deps.go). Это также позволяет легко подменять реализации при необходимости
//
// Почему геттеры не возвращают ошибки:
// Если не удалось подключиться к базе — приложение не может работать. Вместо того,
// чтобы протаскивать ошибку через 5 уровней вызовов, мы логируем и завершаем процесс
// сразу в месте проблемы. Это упрощает API контейнера и код всех вызывающих
type diContainer struct {
	// Инфраструктура
	pgPool *pgxpool.Pool

	// Репозитории
	inventoryRepo part.PartRepository

	// Сервисы
	inventoryService inventoryApi.PartService

	// API-обработчики
	inventoryV1Handler inventoryv1.InventoryServiceServer
}

// PGPool возвращает пул подключений к PostgreSQL
// При первом вызове создаёт пул, проверяет соединение и регистрирует closer
func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		err = pool.Ping(ctx)
		if err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) InventoryRepository(ctx context.Context) part.PartRepository {
	if d.inventoryRepo == nil {
		d.inventoryRepo = inventoryRepository.New(d.PGPool(ctx))
	}

	return d.inventoryRepo
}

func (d *diContainer) InventoryService(ctx context.Context) inventoryApi.PartService {
	if d.inventoryService == nil {
		d.inventoryService = part.New(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryV1Handler == nil {
		d.inventoryV1Handler = inventoryApi.New(d.InventoryService(ctx))
	}

	return d.inventoryV1Handler
}
