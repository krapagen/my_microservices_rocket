package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/krapagen/my_microservices_rocket/inventory/pkg/app"
)

const (
	// Адрес сервера
	grpcAddress = ":50051"

	// gRPC keepalive параметры
	grpcMaxConnectionIdle     = 15 * time.Minute // Закрыть idle-соединения (нет активных RPC)
	grpcMaxConnectionAge      = 30 * time.Minute // Принудительная ротация для балансировки
	grpcMaxConnectionAgeGrace = 5 * time.Second  // Время на завершение активных RPC
	grpcKeepaliveTime         = 5 * time.Minute  // Интервал ping'ов для обнаружения мёртвых соединений
	grpcKeepaliveTimeout      = 1 * time.Second  // Таймаут ожидания pong
	grpcMinPingInterval       = 5 * time.Minute  // Минимальный интервал ping'ов от клиента (защита от DoS)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := godotenv.Load("grpc.env")
	if err != nil {
		slog.Error("ошибка загрузки переменных окружения из grpc.env", "error", err)
		return
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		slog.Error("переменная окружения DB_URI не установлена")
		return
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		slog.Error("ошибка подключения к БД", "error", err)
		return
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("проверка соединения с БД", "error", err)
		return
	}

	slog.Info("подключение к PostgreSQL установлено")

	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		return
	}

	options := make([]grpc.ServerOption, 0, 2+len(app.Interceptors()))
	options = append(
		options,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true, // Разрешить "тёплые" соединения без активных RPC
		}),
	)
	options = append(options, app.Interceptors()...)

	grpcServer := grpc.NewServer(options...)

	app.RegisterServices(grpcServer, pool)

	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	go func() {
		slog.Info("gRPC InventoryService запущен", "address", grpcAddress)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			cancel() // будим main, чтобы не висеть бесконечно
		}
	}()

	<-ctx.Done()
	slog.Info("Остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("Сервер остановлен")
}
