package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/krapagen/my_microservices_rocket/order/pkg/app"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	httpAddress             = ":8080"

	// HTTP таймауты
	httpReadHeaderTimeout = 5 * time.Second  // Защита от Slowloris атаки
	httpReadTimeout       = 15 * time.Second // Лимит на чтение всего запроса
	httpWriteTimeout      = 15 * time.Second // Лимит на запись ответа
	httpIdleTimeout       = 60 * time.Second // Таймаут keep-alive соединений
	httpShutdownTimeout   = 5 * time.Second  // Таймаут graceful shutdown

	// gRPC клиент keepalive параметры
	grpcKeepaliveTime    = 5 * time.Minute // Интервал ping'ов для обнаружения мёртвого сервера
	grpcKeepaliveTimeout = 1 * time.Second // Таймаут ожидания pong
)

func main() {
	// Контекст, который отменяется по SIGINT/SIGTERM или при падении сервера
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Загружаем переменные окружения из .env
	err := godotenv.Load("order.env")
	if err != nil {
		slog.Error("ошибка загрузки переменных окружения из order.env", "error", err)
		return
	}

	// Подключаемся к PostgreSQL
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		slog.Error("переменная окружения DB_URI не установлена")
		return
	}

	// DSN берём из order.env (пока хардкодим в main.go)
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		slog.Error("ошибка подключения к БД", "error", err)
		return
	}
	defer pool.Close()

	// Проверяем соединение
	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("проверка соединения с БД", "error", err)
		return
	}

	slog.Info("подключение к PostgreSQL установлено")

	// 2. Создаём Transaction Manager
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("ошибка создания Transaction Manager", "error", err)
		return
	}

	// Создать gRPC соединение с InventoryService
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return
	}
	defer func() {
		closeErr := inventoryConn.Close()
		if closeErr != nil {
			slog.Error("ошибка закрытия gRPC соединения", "error", closeErr)
		}
	}()

	// Создать gRPC соединение с PaymentService
	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return
	}
	defer func() {
		closeErr := paymentConn.Close()
		if closeErr != nil {
			slog.Error("ошибка закрытия gRPC соединения", "error", closeErr)
		}
	}()

	// Создаём HTTP обработчик с новой архитектурой
	orderHandler, err := app.NewHTTPHandler(
		pool,
		txManager,
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
	)
	if err != nil {
		slog.Error("ошибка создания зависимостей приложения", "error", err)
		return
	}

	// Создаем HTTP сервер с таймаутами для защиты от атак
	// Подробное описание всех параметров: см. week_1/HTTP_SERVER.md
	httpServer := &http.Server{
		Addr:              httpAddress,
		Handler:           orderHandler,
		ReadHeaderTimeout: httpReadHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       httpReadTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      httpWriteTimeout,      // Лимит на запись ответа
		IdleTimeout:       httpIdleTimeout,       // Таймаут keep-alive соединений
	}

	// Запускаем HTTP сервер в горутине
	go func() {
		slog.Info("HTTP Server запущен", "address", httpAddress)
		if serveErr := httpServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска HTTP сервера", "error", serveErr)
			cancel() // будим main, чтобы не висеть бесконечно
		}
	}()

	// Ждём сигнал от ОС или падение сервера
	<-ctx.Done()
	slog.Info("Остановка HTTP сервера")

	// Аккуратно останавливаем HTTP сервер
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancelShutdown()
	if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		slog.Error("ошибка остановки HTTP сервера", "error", shutdownErr)
	}
	slog.Info("HTTP сервер остановлен")
}
