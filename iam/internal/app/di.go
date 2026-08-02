package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	iamauthv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/auth/v1"
	iamuserv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/user/v1"
	"github.com/krapagen/my_microservices_rocket/iam/internal/config"
	iamrepositorysession "github.com/krapagen/my_microservices_rocket/iam/internal/repository/session"
	iamrepositoryuser "github.com/krapagen/my_microservices_rocket/iam/internal/repository/user"
	iamservice "github.com/krapagen/my_microservices_rocket/iam/internal/service/iam"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	platformredis "github.com/krapagen/my_microservices_rocket/platform/pkg/redis"
	iamauthv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
	iamuserv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

type diContainer struct {
	// Инфраструктура
	pgPool      *pgxpool.Pool
	redisClient *redis.Client

	// Менеджер транзакций
	txManager iamrepositoryuser.TxManager

	// Репозитории
	userRepo    iamservice.UserRepository
	sessionRepo iamservice.SessionRepository

	// Сервисы
	iamService iamservice.Service

	// API-обработчики
	iamAuthV1Handler iamauthv1.AuthServiceServer
	iamUserV1Handler iamuserv1.UserServiceServer
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

func (d *diContainer) RedisClient() *redis.Client {
	if d.redisClient == nil {
		client, err := platformredis.NewClient(
			&redis.Options{
				Addr:            config.AppConfig().Redis.Address(),
				DialTimeout:     config.AppConfig().Redis.ConT(),
				ReadTimeout:     config.AppConfig().Redis.ConT(),
				WriteTimeout:    config.AppConfig().Redis.ConT(),
				MaxIdleConns:    config.AppConfig().Redis.MIdle(),
				ConnMaxIdleTime: config.AppConfig().Redis.IdleT(),
				DB:              config.AppConfig().Redis.Database(),
			},
			slog.Default(),
		)
		if err != nil {
			slog.Error("ошибка подключения к Redis", "error", err)
			os.Exit(1)
		}

		closer.Add("redis client", func(_ context.Context) error {
			return client.Close()
		})

		d.redisClient = client
	}

	return d.redisClient
}

func (d *diContainer) TxManager(ctx context.Context) iamrepositoryuser.TxManager {
	if d.txManager == nil {
		txManager, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать Transaction Manager", "error", err)
			os.Exit(1)
		}
		d.txManager = txManager
	}

	return d.txManager
}

func (d *diContainer) IamUserRepository(ctx context.Context) iamservice.UserRepository {
	if d.userRepo == nil {
		d.userRepo = iamrepositoryuser.New(d.PGPool(ctx), d.TxManager(ctx))
	}
	return d.userRepo
}

func (d *diContainer) IamSessionRepository() iamservice.SessionRepository {
	if d.sessionRepo == nil {
		d.sessionRepo = iamrepositorysession.NewRepository(d.RedisClient())
	}
	return d.sessionRepo
}

func (d *diContainer) IamService(ctx context.Context) iamservice.Service {
	if d.iamService == nil {
		d.iamService = iamservice.New(
			d.IamUserRepository(ctx),
			d.IamSessionRepository(),
			config.AppConfig().Session.Duration(),
			bcrypt.DefaultCost,
		)
	}
	return d.iamService
}

func (d *diContainer) IamAuthV1API(ctx context.Context) iamauthv1.AuthServiceServer {
	if d.iamAuthV1Handler == nil {
		d.iamAuthV1Handler = iamauthv1api.New(d.IamService(ctx))
	}
	return d.iamAuthV1Handler
}

func (d *diContainer) IamUserV1API(ctx context.Context) iamuserv1.UserServiceServer {
	if d.iamUserV1Handler == nil {
		d.iamUserV1Handler = iamuserv1api.New(d.IamService(ctx))
	}
	return d.iamUserV1Handler
}
