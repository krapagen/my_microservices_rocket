package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	op := "order/internal/service/order/Cancel"
	log := slog.With("op", op)
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		log.ErrorContext(ctx, "не удалось получить заказ", "orderUUID", orderUUID, "error", err)
		if errors.Is(err, errs.ErrOrderNotFound) {
			log.ErrorContext(ctx, "Деталь не найдена в хранилище")
			return errs.ErrOrderNotFound
		}
		return fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		log.ErrorContext(ctx, "Заказ уже оплачен", "orderUUID", orderUUID)
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		log.ErrorContext(ctx, "Заказ уже отменен", "orderUUID", orderUUID)
		return errs.ErrOrderCancelled
	case model.OrderStatusPendingPayment:
		log.InfoContext(ctx, "Отмена заказа", "orderUUID", orderUUID)
	default:
		log.ErrorContext(ctx, "Неизвестный статус заказа", "orderUUID", orderUUID)
		return errs.ErrOrderStatusIncorrect
	}

	order.Status = model.OrderStatusCancelled

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		log.ErrorContext(ctx, "не удалось отменить заказ", "orderUUID", orderUUID, "error", err)
		return fmt.Errorf("обновить заказ: %w", err)
	}

	if err := s.inventoryClient.ReleaseParts(ctx, partUUIDsFromOrderItems(order.Items)); err != nil {
		log.ErrorContext(ctx, "не удалось освободить детали при отмене заказа", "orderUUID", orderUUID, "error", err)
		return fmt.Errorf("освободить детали: %w", err)
	}

	log.InfoContext(ctx, "Заказ успешно отменен", "orderUUID", orderUUID)
	return nil
}
