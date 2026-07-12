package v1

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	op := "order/internal/api/order/v1/CancelOrder"
	log := slog.With("op", op)
	if params.OrderUUID == uuid.Nil {
		log.ErrorContext(ctx, "Номер заказа неверный", "error", errs.ErrInvalidUUID)
		return nil, errs.ErrInvalidUUID
	}
	log.InfoContext(ctx, "Отмена заказа...", "order_uuid", params.OrderUUID)
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		log.ErrorContext(ctx, "Ошибка отмены заказа", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "Заказ успешно отменен", "order_uuid", params.OrderUUID)
	return &orderv1.CancelOrderResponse{}, nil
}
