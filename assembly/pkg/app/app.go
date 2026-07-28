package app

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/consumer/order_paid"
	"github.com/krapagen/my_microservices_rocket/assembly/internal/producer/ship_assembled"
	"github.com/krapagen/my_microservices_rocket/assembly/internal/service/assembly"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka/consumer"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka/producer"
)

// Config holds test configuration for the assembly application.
type Config struct {
	OrderPaidTopic     string
	ShipAssembledTopic string
	MinBuildTimeSec    int64
	MaxBuildTimeSec    int64
}

// App is a test assembly application wiring consumer and producer.
type App struct {
	svc order_paid.ConsumerService
}

// New creates a test assembly app from provided Kafka clients and config.
// It uses the same internal services as the production app.
func New(syncProducer sarama.SyncProducer, consumerGroup sarama.ConsumerGroup, cfg Config) *App {
	if cfg.MaxBuildTimeSec < 0 {
		cfg.MaxBuildTimeSec = 0
	}

	buildTime := time.Duration(cfg.MaxBuildTimeSec) * time.Second

	orderPaidConsumer := consumer.NewConsumer(
		consumerGroup,
		[]string{cfg.OrderPaidTopic},
	)
	shipAssembledProducer := producer.NewProducer(
		syncProducer,
		cfg.ShipAssembledTopic,
	)

	assembler := assembly.NewService(buildTime)
	shipAssembledSvc := ship_assembled.NewService(shipAssembledProducer)
	svc := order_paid.NewService(orderPaidConsumer, assembler, shipAssembledSvc)

	return &App{svc: svc}
}

// RunConsumer starts the OrderPaid consumer.
func (a *App) RunConsumer(ctx context.Context) error {
	return a.svc.RunConsumer(ctx)
}
