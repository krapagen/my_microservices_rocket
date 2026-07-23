package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	slots := model.ShipSlots{
		HullUUID:   in.HullUUID,
		EngineUUID: in.EngineUUID,
		ShieldUUID: in.ShieldUUID,
		WeaponUUID: in.WeaponUUID,
	}

	if err := s.inventoryClient.ValidateCompatibility(ctx, slots); err != nil {
		log.ErrorContext(ctx, "детали несовместимы", "error", err)
		return model.Order{}, fmt.Errorf("проверить совместимость: %w", mapInventoryError(err))
	}

	partUUIDs := in.PartUUIDs()
	if err := s.inventoryClient.ReserveParts(ctx, partUUIDs); err != nil {
		log.ErrorContext(ctx, "не удалось зарезервировать детали", "error", err)
		return model.Order{}, fmt.Errorf("зарезервировать детали: %w", mapInventoryError(err))
	}

	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		if releaseErr := s.inventoryClient.ReleaseParts(ctx, partUUIDs); releaseErr != nil {
			log.ErrorContext(ctx, "не удалось освободить детали после ошибки получения деталей", "error", releaseErr)
		}
		log.ErrorContext(ctx, "не удалось получить детали из инвентаря", "error", err)
		return model.Order{}, fmt.Errorf("получить детали: %w", mapInventoryError(err))
	}
	if len(parts) != len(partUUIDs) {
		if releaseErr := s.inventoryClient.ReleaseParts(ctx, partUUIDs); releaseErr != nil {
			log.ErrorContext(ctx, "не удалось освободить детали при неполном списке деталей", "error", releaseErr)
		}
		log.ErrorContext(ctx, "не все детали найдены в инвентаре", "expected", len(partUUIDs), "actual", len(parts))
		return model.Order{}, fmt.Errorf("получить детали: %w", errs.ErrPartNotFound)
	}
	log.InfoContext(ctx, "Детали получены из инвентаря", "parts", parts)

	items := make([]model.OrderItem, 0, len(parts))
	for _, part := range parts {
		items = append(items, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
		log.InfoContext(ctx, "Деталь добавлена в заказ", "part.Name", part.Name, "part.UUID", part.UUID, "part.PartType", part.PartType, "part.Price", part.Price)
	}

	order := model.Order{
		UUID:      uuid.New(),
		Items:     items,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		if releaseErr := s.inventoryClient.ReleaseParts(ctx, partUUIDs); releaseErr != nil {
			log.ErrorContext(ctx, "не удалось освободить детали после ошибки создания заказа", "error", releaseErr)
		}
		log.ErrorContext(ctx, "не удалось создать заказ", "error", err)
		return model.Order{}, fmt.Errorf("создать заказ: %w", err)
	}
	log.InfoContext(ctx, "Заказ успешно создан", "order.UUID", order.UUID, "order.Status", order.Status, "order.CreatedAt", order.CreatedAt)
	return order, nil
}

func mapInventoryError(err error) error {
	switch status.Code(err) {
	case codes.NotFound:
		return errs.ErrPartNotFound
	case codes.InvalidArgument:
		return errs.ErrPartTypeMismatch
	case codes.FailedPrecondition:
		return errs.ErrIncompatibleParts
	case codes.ResourceExhausted:
		return errs.ErrOutOfStock
	default:
		return err
	}
}
