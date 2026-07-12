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

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	op := "order/internal/service/order/Pay"
	log := slog.With("op", op)
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		log.ErrorContext(ctx, "Не удалось получить заказ", "orderUUID", orderUUID, "error", err)
		if errors.Is(err, errs.ErrOrderNotFound) {
			log.ErrorContext(ctx, "Заказ не найден", "orderUUID", orderUUID)
			return uuid.Nil, errs.ErrOrderNotFound
		}
		return uuid.Nil, fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		log.ErrorContext(ctx, "Заказ уже оплачен", "orderUUID", orderUUID)
		return uuid.Nil, errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		log.ErrorContext(ctx, "Заказ отменен", "orderUUID", orderUUID)
		return uuid.Nil, errs.ErrOrderCancelled
	case model.OrderStatusPendingPayment:
		log.InfoContext(ctx, "Оплата заказа", "orderUUID", orderUUID, "method", method)
	default:
		log.ErrorContext(ctx, "Неизвестный статус заказа", "orderUUID", orderUUID, "status", order.Status)
		return uuid.Nil, errs.ErrOrderStatusIncorrect
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, method)
	if err != nil {
		log.ErrorContext(ctx, "Не удалось оплатить заказ", "orderUUID", orderUUID, "method", method, "error", err)
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}
	log.InfoContext(ctx, "Заказ успешно оплачен", "orderUUID", orderUUID, "method", method, "transactionUUID", transactionUUID)
	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		log.ErrorContext(ctx, "Не удалось обновить заказ после оплаты", "orderUUID", orderUUID, "error", err)
		return uuid.Nil, fmt.Errorf("обновить заказ: %w", err)
	}
	log.InfoContext(ctx, "Заказ успешно обновлен после оплаты", "orderUUID", orderUUID, "transactionUUID", transactionUUID)
	return transactionUUID, nil
}
