package assembly_consumer

import (
	"context"
	"log/slog"
)

type service struct {
	consumer        Consumer
	orderRepository OrderRepository
	inventoryClient InventoryClient
	txManager       TxManager
}

func NewService(
	consumer Consumer,
	orderRepository OrderRepository,
	inventoryClient InventoryClient,
	txManager TxManager,
) ConsumerService {
	return &service{
		consumer:        consumer,
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		txManager:       txManager,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя ShipAssembled")
	return s.consumer.Consume(ctx, s.ShipAssembledHandler)
}
