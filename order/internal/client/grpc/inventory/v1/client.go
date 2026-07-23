package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

// New создаёт обёртку над gRPC клиентом InventoryService.
func New(c inventoryv1.InventoryServiceClient) *client {
	return &client{inventoryClient: c}
}

func (c *client) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	op := "order/internal/client/grpc/inventory/v1/ListParts"
	log := slog.With("op", op)
	uuidsStrings := make([]string, 0, len(uuids))
	for _, uuidCur := range uuids {
		uuidsStrings = append(uuidsStrings, uuidCur.String())
	}
	resp, err := c.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuidsStrings,
	})
	if err != nil {
		log.ErrorContext(ctx, "Не удалось получить детали из InventoryService", "error", err)
		switch status.Code(err) {
		case codes.NotFound:
			log.ErrorContext(ctx, "Детали не найдены", "uuids", uuidsStrings)
			return nil, errs.ErrPartNotFound
		default:
			log.ErrorContext(ctx, "Ошибка при вызове ListParts", "error", err)
			return nil, fmt.Errorf("получить список деталей: %w", err)
		}
	}
	log.InfoContext(ctx, "Детали успешно получены из InventoryService", "count", len(resp.GetParts()))
	return converter.PartsToModel(resp.GetParts()), nil
}

func (c *client) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	op := "order/internal/client/grpc/inventory/v1/ValidateCompatibility"
	log := slog.With("op", op)

	_, err := c.inventoryClient.ValidateCompatibility(ctx, converter.ShipSlotsToProto(slots))
	if err != nil {
		log.ErrorContext(ctx, "Ошибка проверки совместимости в InventoryService", "error", err)
		return fmt.Errorf("проверить совместимость: %w", err)
	}

	return nil
}

func (c *client) ReserveParts(ctx context.Context, uuids []uuid.UUID) error {
	op := "order/internal/client/grpc/inventory/v1/ReserveParts"
	log := slog.With("op", op)

	uuidsStrings := make([]string, 0, len(uuids))
	for _, uuidCur := range uuids {
		uuidsStrings = append(uuidsStrings, uuidCur.String())
	}

	_, err := c.inventoryClient.ReserveParts(ctx, &inventoryv1.ReservePartsRequest{Uuids: uuidsStrings})
	if err != nil {
		log.ErrorContext(ctx, "Ошибка резервирования деталей в InventoryService", "error", err)
		return fmt.Errorf("зарезервировать детали: %w", err)
	}

	return nil
}

func (c *client) ReleaseParts(ctx context.Context, uuids []uuid.UUID) error {
	op := "order/internal/client/grpc/inventory/v1/ReleaseParts"
	log := slog.With("op", op)

	uuidsStrings := make([]string, 0, len(uuids))
	for _, uuidCur := range uuids {
		uuidsStrings = append(uuidsStrings, uuidCur.String())
	}

	_, err := c.inventoryClient.ReleaseParts(ctx, &inventoryv1.ReleasePartsRequest{Uuids: uuidsStrings})
	if err != nil {
		log.ErrorContext(ctx, "Ошибка освобождения деталей в InventoryService", "error", err)
		return fmt.Errorf("освободить детали: %w", err)
	}

	return nil
}
