package assembly_consumer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
)

func (s *service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := decodeShipAssembled(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось декодировать ShipAssembled", "error", err)
		return err
	}

	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepository.GetForUpdate(txCtx, event.OrderUUID)
		if err != nil {
			slog.ErrorContext(txCtx, "не удалось получить заказ для обновления", "orderUUID", event.OrderUUID, "error", err)
			return fmt.Errorf("получить заказ: %w", err)
		}

		switch order.Status {
		case model.OrderStatusPaid:
			// ожидаемый статус перед переводом заказа в ASSEMBLED
		case model.OrderStatusAssembled:
			// идемпотентность для повторной доставки события
			slog.InfoContext(txCtx, "заказ уже собран, повторное сообщение пропущено", "orderUUID", order.UUID)
			return nil
		case model.OrderStatusCancelled:
			return errs.ErrOrderCancelled
		case model.OrderStatusPendingPayment:
			return errs.ErrOrderStatusIncorrect
		default:
			return errs.ErrOrderStatusIncorrect
		}

		partUUIDs := partUUIDsFromOrderItems(order.Items)
		if err = s.inventoryClient.CommitParts(txCtx, partUUIDs); err != nil {
			slog.ErrorContext(txCtx, "не удалось списать детали в inventory", "orderUUID", order.UUID, "error", err)
			return fmt.Errorf("commit parts: %w", err)
		}

		order.Status = model.OrderStatusAssembled
		if err = s.orderRepository.Update(txCtx, order); err != nil {
			slog.ErrorContext(txCtx, "не удалось обновить статус заказа", "orderUUID", order.UUID, "error", err)
			return fmt.Errorf("обновить заказ: %w", err)
		}

		slog.InfoContext(txCtx, "заказ переведен в ASSEMBLED", "orderUUID", order.UUID)
		return nil
	})
}

func partUUIDsFromOrderItems(items []model.OrderItem) []uuid.UUID {
	uuids := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		uuids = append(uuids, item.PartUUID)
	}
	return uuids
}
