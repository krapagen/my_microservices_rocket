package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/input"
)

func (s *service) Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error) {
	op := "order/internal/service/order/Create"
	log := slog.With("op", op)
	if (in.HullUUID == uuid.Nil) || (in.EngineUUID == uuid.Nil) {
		log.ErrorContext(ctx, "не указаны обязательные детали", "HullUUID", in.HullUUID, "EngineUUID", in.EngineUUID)
		return model.Order{}, fmt.Errorf("не указаны обязательные детали: %w", errs.ErrMissingRequiredParts)
	}
	log.InfoContext(ctx, "Создание заказа", "HullUUID", in.HullUUID, "EngineUUID", in.EngineUUID, "ShieldUUID", in.ShieldUUID, "WeaponUUID", in.WeaponUUID)
	parts, err := s.inventoryClient.ListParts(ctx, in.PartUUIDs())
	if err != nil {
		log.ErrorContext(ctx, "не удалось получить детали из инвентаря", "error", err)
		return model.Order{}, fmt.Errorf("получить детали: %w", err)
	}
	log.InfoContext(ctx, "Детали получены из инвентаря", "parts", parts)
	items := make([]model.OrderItem, 0, len(parts))
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			log.ErrorContext(ctx, "деталь отсутствует на складе", "part.Name", part.Name, "part.UUID", part.UUID)
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		items = append(items, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
		log.InfoContext(ctx, "Деталь добавлена в заказ", "part.Name", part.Name, "part.UUID", part.UUID, "part.PartType", part.PartType, "part.Price", part.Price)
	}

	// Проверяем, что в заказе есть детали
	if len(items) == 0 {
		log.ErrorContext(ctx, "заказ не содержит деталей")
		return model.Order{}, fmt.Errorf("заказ должен содержать хотя бы одну деталь: %w", errs.ErrMissingRequiredParts)
	}

	order := model.Order{
		UUID:      uuid.New(),
		Items:     items,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}
	err = s.orderRepo.Create(ctx, order)
	if err != nil {

		log.ErrorContext(ctx, "не удалось создать заказ", "error", err)
		return model.Order{}, fmt.Errorf("создать заказ: %w", err)
	}
	log.InfoContext(ctx, "Заказ успешно создан", "order.UUID", order.UUID, "order.Status", order.Status, "order.CreatedAt", order.CreatedAt)
	return order, nil
}
