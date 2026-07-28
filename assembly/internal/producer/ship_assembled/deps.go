package ship_assembled

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
)

// Producer определяет контракт для отправки сообщений в Kafka
type Producer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}

// ProducerService определяет контракт публикации ShipAssembled.
type ProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}
