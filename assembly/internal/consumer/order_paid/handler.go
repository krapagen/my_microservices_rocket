package order_paid

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeOrderPaid(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать OrderPaid", "error", err)
		return err
	}

	slog.InfoContext(
		ctx, "обработка сообщения",
		"topic", msg.Topic,
		"partition", msg.Partition,
		"offset", msg.Offset,
		"event_uuid", event.EventUUID,
		"order_uuid", event.OrderUUID,
		"user_uuid", event.UserUUID,
	)

	assembled, err := s.assembler.Assemble(ctx, event)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось собрать корабль", "error", err)
		return fmt.Errorf("assemble ship: %w", err)
	}

	if err = s.shipAssembledProducer.ProduceShipAssembled(ctx, assembled); err != nil {
		slog.ErrorContext(ctx, "не удалось отправить ShipAssembled", "error", err)
		return fmt.Errorf("produce ship assembled: %w", err)
	}

	return nil
}
