package v1

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/api/converter"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	op := "order/internal/api/order/v1/PayOrder"
	log := slog.With("op", op)
	if params.OrderUUID == uuid.Nil {
		log.ErrorContext(ctx, "Номер заказа неверный", "error", errs.ErrInvalidUUID)
		return nil, errs.ErrInvalidUUID
	}

	// Конвертация запроса
	payParams := converter.PayOrderRequestToModel(req, params)

	// Вызов сервиса
	log.InfoContext(ctx, "Оплата заказа...", "order_uuid", payParams.OrderUUID, "payment_method", payParams.PaymentMethod)
	transactionUUID, err := a.orderService.Pay(ctx, payParams.OrderUUID, payParams.PaymentMethod)
	if err != nil {
		log.ErrorContext(ctx, "Оплата заказа не прошла", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "Оплата заказа прошла успешно", "order_uuid", payParams.OrderUUID, "transaction_uuid", transactionUUID)
	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
