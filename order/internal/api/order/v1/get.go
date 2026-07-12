package v1

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/api/converter"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	op := "order/internal/api/order/v1/GetOrder"
	log := slog.With("op", op)
	if params.OrderUUID == uuid.Nil {
		log.ErrorContext(ctx, "Номер заказа неверный", "error", errs.ErrInvalidUUID)
		return nil, errs.ErrInvalidUUID
	}
	log.InfoContext(ctx, "Получение заказа...", "order_uuid", params.OrderUUID)
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		log.ErrorContext(ctx, "Ошибка получения заказа", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "Заказ успешно получен", "order_uuid", params.OrderUUID)
	return converter.OrderModelToDTO(order), nil
}
