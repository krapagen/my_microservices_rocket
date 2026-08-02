package app

import (
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	iamauthv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/auth/v1"
	iamuserv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/user/v1"
	"github.com/krapagen/my_microservices_rocket/iam/internal/interceptor"
	iamrepositorysession "github.com/krapagen/my_microservices_rocket/iam/internal/repository/session"
	iamrepositoryuser "github.com/krapagen/my_microservices_rocket/iam/internal/repository/user"
	iamservice "github.com/krapagen/my_microservices_rocket/iam/internal/service/iam"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

// NewGRPCServer создаёт и регистрирует gRPC-сервер IAM.
func NewGRPCServer(pool *pgxpool.Pool, rdb *redis.Client, sessionTTL time.Duration, bcryptCost int) *grpc.Server {
	s := grpc.NewServer(Interceptors()...)
	RegisterServices(s, pool, rdb, sessionTTL, bcryptCost)
	return s
}

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool, rdb *redis.Client, sessionTTL time.Duration, bcryptCost int, txManagers ...iamrepositoryuser.TxManager) {
	var txManager iamrepositoryuser.TxManager
	if len(txManagers) > 0 {
		txManager = txManagers[0]
	} else {
		var err error
		txManager, err = manager.New(trmpgx.NewDefaultFactory(pool))
		if err != nil {
			panic(fmt.Errorf("create transaction manager: %w", err))
		}
	}

	userRepo := iamrepositoryuser.New(pool, txManager)
	sessionRepo := iamrepositorysession.NewRepository(rdb)
	svc := iamservice.New(userRepo, sessionRepo, sessionTTL, bcryptCost)

	authAPI := iamauthv1api.New(svc)
	userAPI := iamuserv1api.New(svc)

	authv1.RegisterAuthServiceServer(grpcServer, authAPI)
	userv1.RegisterUserServiceServer(grpcServer, userAPI)
}
