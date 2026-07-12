package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	options := make([]grpc.ServerOption, 0, 2+len(app.Interceptors()))
	options = append(options,
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

	app.RegisterServices(grpcServer)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("gRPC InventoryService запущен", "address", grpcAddress)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			cancel() // будим main, чтобы не висеть бесконечно
		}
	}()

	// Ждём сигнал от ОС или падение сервера
	<-ctx.Done()
	slog.Info("Остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("Сервер остановлен")
}
