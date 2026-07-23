package order

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	op := "order/internal/service/order/Pay"
	log := slog.With("op", op)
	var transactionUUID uuid.UUID

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. Читаем заказ в транзакции
		order, err := s.orderRepo.Get(txCtx, orderUUID)
		if err != nil {
			log.ErrorContext(txCtx, "Не удалось получить заказ", "orderUUID", orderUUID, "error", err)
			return fmt.Errorf("получить заказ: %w", err)
		}
		log.InfoContext(txCtx, "Заказ успешно прочитан", "orderUUID", orderUUID)
		// 2. Проверяем статус
		switch order.Status {
		case model.OrderStatusPendingPayment:
		case model.OrderStatusPaid:
			log.ErrorContext(txCtx, "Заказ уже оплачен", "orderUUID", orderUUID)
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			log.ErrorContext(txCtx, "Заказ отменён", "orderUUID", orderUUID)
			return errs.ErrOrderCancelled
		default:
			log.ErrorContext(txCtx, "Неизвестный статус заказа", "orderUUID", orderUUID, "status", order.Status)
			return errs.ErrOrderStatusIncorrect
		}

		// 3. Вызываем PaymentService (gRPC внутри транзакции — учебный пример)
		transactionUUID, err = s.paymentClient.PayOrder(txCtx, orderUUID, method)
		if err != nil {
			log.ErrorContext(txCtx, "Не удалось оплатить заказ", "orderUUID", orderUUID, "method", method, "error", err)
			return fmt.Errorf("оплатить заказ: %w", err)
		}
		log.InfoContext(txCtx, "Заказ успешно оплачен", "orderUUID", orderUUID, "method", method, "transactionUUID", transactionUUID)
		// 4. Обновляем заказ
		order.Status = model.OrderStatusPaid
		order.TransactionUUID = &transactionUUID
		order.PaymentMethod = &method

		err = s.orderRepo.Update(txCtx, order)
		if err != nil {
			log.ErrorContext(txCtx, "Не удалось обновить заказ после оплаты", "orderUUID", orderUUID, "error", err)
			return fmt.Errorf("обновить заказ: %w", err)
		}
		log.InfoContext(txCtx, "Заказ успешно обновлён после оплаты", "orderUUID", orderUUID)
		return nil
	})
	if err != nil {
		log.ErrorContext(ctx, "Ошибка при оплате заказа", "orderUUID", orderUUID, "method", method, "error", err)
		return uuid.Nil, err
	}
	log.InfoContext(ctx, "Оплата заказа успешно завершена", "orderUUID", orderUUID, "method", method, "transactionUUID", transactionUUID)
	return transactionUUID, nil
}
