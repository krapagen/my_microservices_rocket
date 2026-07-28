package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/IBM/sarama"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	orderApi "github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	inventoryClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/payment/v1"
	"github.com/krapagen/my_microservices_rocket/order/internal/config"
	assemblyConsumer "github.com/krapagen/my_microservices_rocket/order/internal/consumer/assembly_consumer"
	orderProducer "github.com/krapagen/my_microservices_rocket/order/internal/producer/order_producer"
	orderRepository "github.com/krapagen/my_microservices_rocket/order/internal/repository/order"
	orderService "github.com/krapagen/my_microservices_rocket/order/internal/service/order"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/krapagen/my_microservices_rocket/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/krapagen/my_microservices_rocket/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/krapagen/my_microservices_rocket/platform/pkg/middleware/kafka"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	// Pool подключений
	pgPool *pgxpool.Pool

	// Менеджер транзакций
	txManager orderRepository.TxManager

	// Репозитории
	repo orderService.OrderRepository

	// inventoryClient
	inventoryConn   *grpc.ClientConn
	inventoryClient orderService.InventoryClient

	// paymentClient
	paymentConn   *grpc.ClientConn
	paymentClient orderService.PaymentClient

	// kafka producer
	syncProducer      sarama.SyncProducer
	consumerGroup     sarama.ConsumerGroup
	orderPaidProducer *wrappedKafkaProducer.Producer
	shipAssembledCons *wrappedKafkaConsumer.Consumer

	// producer service
	orderPaidProducerService orderService.OrderPaidProducer
	assemblyConsumerService  assemblyConsumer.ConsumerService

	// Сервисы
	svc orderApi.OrderService

	// API-обработчики
	handler orderv1.Handler
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

func (d *diContainer) TxManager(ctx context.Context) orderRepository.TxManager {
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

// OrderRepository возвращает репозиторий
func (d *diContainer) OrderRepository(ctx context.Context) orderService.OrderRepository {
	if d.repo == nil {
		d.repo = orderRepository.New(d.PGPool(ctx), d.TxManager(ctx))
	}

	return d.repo
}

func (d *diContainer) InventoryClient() orderService.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryClient.New(inventoryv1.NewInventoryServiceClient(d.inventoryConn))
		closer.Add("Inventory gRPC client", func(_ context.Context) error {
			return d.inventoryConn.Close()
		})
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentClient() orderService.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentClient.New(paymentv1.NewPaymentServiceClient(d.paymentConn))
		closer.Add("Payment gRPC client", func(_ context.Context) error {
			return d.paymentConn.Close()
		})
	}
	return d.paymentClient
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		cfg.Producer.RequiredAcks = sarama.WaitForAll

		p, err := sarama.NewSyncProducer(config.AppConfig().Kafka.KafkaBrokers(), cfg)
		if err != nil {
			slog.Error("не удалось создать Kafka sync producer", "error", err)
			os.Exit(1)
		}
		closer.Add("Kafka sync producer", func(_ context.Context) error {
			return p.Close()
		})
		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
			sarama.NewBalanceStrategyRoundRobin(),
		}

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.KafkaBrokers(),
			config.AppConfig().ShipAssembledConsumer.ShipAssembledGroup(),
			cfg,
		)
		if err != nil {
			slog.Error("не удалось создать Kafka consumer group", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return group.Close()
		})
		d.consumerGroup = group
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidProducer() *wrappedKafkaProducer.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.OrderPaidProducerTopic(),
		)
	}

	return d.orderPaidProducer
}

func (d *diContainer) OrderPaidProducerService() orderService.OrderPaidProducer {
	if d.orderPaidProducerService == nil {
		d.orderPaidProducerService = orderProducer.New(d.OrderPaidProducer())
	}

	return d.orderPaidProducerService
}

func (d *diContainer) ShipAssembledConsumer() *wrappedKafkaConsumer.Consumer {
	if d.shipAssembledCons == nil {
		d.shipAssembledCons = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().ShipAssembledConsumer.ShipAssembledTopicName()},
			wrappedKafkaConsumer.WithMiddlewares(kafkaMiddleware.ConsumerLogging()),
		)
	}

	return d.shipAssembledCons
}

func (d *diContainer) AssemblyConsumerService(ctx context.Context) assemblyConsumer.ConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = assemblyConsumer.NewService(
			d.ShipAssembledConsumer(),
			d.OrderRepository(ctx),
			d.InventoryClient(),
			d.TxManager(ctx),
		)
	}

	return d.assemblyConsumerService
}

// OrderService возвращает сервис бизнес-логики заказов
func (d *diContainer) OrderService(ctx context.Context) orderApi.OrderService {
	if d.svc == nil {
		d.svc = orderService.New(d.OrderRepository(ctx),
			d.InventoryClient(),
			d.PaymentClient(),
			d.OrderPaidProducerService(),
			d.TxManager(ctx),
		)
	}
	return d.svc
}

// OrderV1API возвращает обработчик сервиса заказов
func (d *diContainer) OrderV1API(ctx context.Context) orderv1.Handler {
	if d.handler == nil {
		d.handler = orderApi.NewAPI(d.OrderService(ctx))
	}

	return d.handler
}
