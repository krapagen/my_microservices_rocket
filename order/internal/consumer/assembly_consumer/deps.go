package assembly_consumer

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
)

// Consumer определяет контракт для потребления сообщений из Kafka.
type Consumer interface {
	Consume(ctx context.Context, handler kafka.MessageHandler) error
}

// OrderRepository определяет операции хранения заказа для обработчика ShipAssembled.
type OrderRepository interface {
	GetForUpdate(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

// InventoryClient определяет контракт списания деталей в InventoryService.
type InventoryClient interface {
	CommitParts(ctx context.Context, uuids []uuid.UUID) error
}

// TxManager определяет контракт выполнения бизнес-логики в транзакции.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// ConsumerService определяет контракт запуска Kafka consumer.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
