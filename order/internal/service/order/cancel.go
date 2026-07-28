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
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.GetForUpdate(txCtx, orderUUID)
		if err != nil {
			log.ErrorContext(txCtx, "не удалось получить заказ", "orderUUID", orderUUID, "error", err)
			if errors.Is(err, errs.ErrOrderNotFound) {
				log.ErrorContext(txCtx, "Заказ не найден в хранилище", "orderUUID", orderUUID)
				return errs.ErrOrderNotFound
			}
			return fmt.Errorf("получить заказ: %w", err)
		}

		switch order.Status {
		case model.OrderStatusPaid:
			log.ErrorContext(txCtx, "Заказ уже оплачен", "orderUUID", orderUUID)
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			log.ErrorContext(txCtx, "Заказ уже отменен", "orderUUID", orderUUID)
			return errs.ErrOrderCancelled
		case model.OrderStatusAssembled:
			log.ErrorContext(txCtx, "Заказ уже собран", "orderUUID", orderUUID)
			return errs.ErrOrderAssembled
		case model.OrderStatusPendingPayment:
			log.InfoContext(txCtx, "Отмена заказа", "orderUUID", orderUUID)
		default:
			log.ErrorContext(txCtx, "Неизвестный статус заказа", "orderUUID", orderUUID)
			return errs.ErrOrderStatusIncorrect
		}

		order.Status = model.OrderStatusCancelled

		err = s.orderRepo.Update(txCtx, order)
		if err != nil {
			log.ErrorContext(txCtx, "не удалось отменить заказ", "orderUUID", orderUUID, "error", err)
			return fmt.Errorf("обновить заказ: %w", err)
		}

		if err := s.inventoryClient.ReleaseParts(txCtx, partUUIDsFromOrderItems(order.Items)); err != nil {
			log.ErrorContext(txCtx, "не удалось освободить детали при отмене заказа", "orderUUID", orderUUID, "error", err)
			return fmt.Errorf("освободить детали: %w", err)
		}
		return nil
	})
	if err != nil {
		log.ErrorContext(ctx, "Ошибка при отмене заказа", "orderUUID", orderUUID, "error", err)
		return err
	}
	log.InfoContext(ctx, "Заказ успешно отменен", "orderUUID", orderUUID)
	return nil
}
