package order_paid

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
)

// Consumer определяет контракт для потребления сообщений из Kafka
type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

// Assembler определяет контракт бизнес-логики сборки корабля.
type Assembler interface {
	Assemble(ctx context.Context, in model.OrderPaid) (model.ShipAssembled, error)
}

// ShipAssembledProducer определяет контракт публикации события ShipAssembled.
type ShipAssembledProducer interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}

// ConsumerService определяет контракт запуска Kafka consumer.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
