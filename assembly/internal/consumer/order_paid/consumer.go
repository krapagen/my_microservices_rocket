package order_paid

import (
	"context"
	"log/slog"
)

type service struct {
	orderPaidConsumer     Consumer
	assembler             Assembler
	shipAssembledProducer ShipAssembledProducer
}

func NewService(
	orderPaidConsumer Consumer,
	assembler Assembler,
	shipAssembledProducer ShipAssembledProducer,
) ConsumerService {
	return &service{
		orderPaidConsumer:     orderPaidConsumer,
		assembler:             assembler,
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "запуск потребителя OrderPaid")

	return s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
}
