package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapi "github.com/krapagen/my_microservices_rocket/inventory/internal/api/inventory/v1"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/interceptor"
	repository "github.com/krapagen/my_microservices_rocket/inventory/internal/repository/part"
	service "github.com/krapagen/my_microservices_rocket/inventory/internal/service/part"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

// Interceptors возвращает grpc.ServerOption для тестов
func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

// RegisterServices регистрирует сервисы на gRPC сервере
func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	repo := repository.New(pool)
	svc := service.New(repo)
	api := inventoryapi.New(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}
