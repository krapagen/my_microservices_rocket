package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/converter"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ReleaseParts(
	ctx context.Context,
	req *inventoryv1.ReleasePartsRequest,
) (*inventoryv1.ReleasePartsResponse, error) {
	op := "Функция inventory/internal/api/inventory/v1/ReleaseParts"
	log := slog.With("op", op)

	convert := converter.NewConverter()
	uuids, err := convert.ToGetInputs(req.GetUuids())
	if err != nil {
		log.ErrorContext(ctx, "неверный формат uuid", "error", err)
		return nil, err
	}

	if err := a.partService.Release(ctx, uuids); err != nil {
		log.ErrorContext(ctx, "ошибка освобождения деталей", "error", err)
		return nil, err
	}

	return &inventoryv1.ReleasePartsResponse{}, nil
}
