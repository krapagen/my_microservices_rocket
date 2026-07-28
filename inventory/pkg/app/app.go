package app

import (
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	inventoryapi "github.com/krapagen/my_microservices_rocket/inventory/internal/api/inventory/v1"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/interceptor"
	repository "github.com/krapagen/my_microservices_rocket/inventory/internal/repository/part"
	service "github.com/krapagen/my_microservices_rocket/inventory/internal/service/application/part"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

// Interceptors возвращает grpc.ServerOption для тестов
func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

// RegisterServices регистрирует сервисы на gRPC сервере.
// txManager передаётся опционально: если не передан — создаётся из pool.
func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool, txManagers ...repository.TxManager) {
	var txManager repository.TxManager
	if len(txManagers) > 0 {
		txManager = txManagers[0]
	} else {
		var err error
		txManager, err = manager.New(trmpgx.NewDefaultFactory(pool))
		if err != nil {
			panic(fmt.Errorf("create transaction manager: %w", err))
		}
	}

	repo := repository.New(pool, txManager)
	svc := service.New(repo, txManager)
	api := inventoryapi.New(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}
