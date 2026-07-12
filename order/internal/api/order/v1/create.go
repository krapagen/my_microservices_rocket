package v1

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/api/converter"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	op := "order/internal/api/order/v1/CreateOrder"
	log := slog.With("op", op)
	if req.HullUUID == uuid.Nil || req.EngineUUID == uuid.Nil {
		log.ErrorContext(ctx, "Отсутствуют обязательные части заказа", "error", errs.ErrMissingRequiredParts)
		return nil, errs.ErrMissingRequiredParts
	}
	log.InfoContext(ctx, "Создание заказа...", "hull_uuid", req.HullUUID, "engine_uuid", req.EngineUUID)
	reqInput := converter.CreateOrderRequestToInput(req)
	respModel, err := a.orderService.Create(ctx, reqInput)
	if err != nil {
		log.ErrorContext(ctx, "Не получилось создать заказ", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "Заказ успешно создан", "order_uuid", respModel.UUID)
	return &orderv1.CreateOrderResponse{
		OrderUUID:  respModel.UUID,
		TotalPrice: respModel.TotalPrice(),
	}, nil
}
