package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	op := "Функция inventory/internal/api/inventory/v1/ListParts"
	log := slog.With("op", op)
	rawUuids := req.GetUuids()
	rawPartType := req.GetPartType()
	log.InfoContext(ctx, "получен запрос на список деталей", "uuids", rawUuids, "partType", rawPartType)
	convert := converter.NewConverter()
	uuids, err := convert.ToGetInputs(rawUuids)
	if err != nil {
		log.ErrorContext(ctx, "неверный формат uuid", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "валидный формат uuid", "uuids", uuids)
	partType := convert.DtoTypeToPartType(rawPartType)

	models, err := a.partService.List(ctx, input.PartFilter{UUIDs: uuids, PartType: partType})
	if err != nil {
		log.ErrorContext(ctx, "ошибка при получении списка деталей", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "список деталей получен", "count", len(models))
	return &inventoryv1.ListPartsResponse{
		Parts: convert.PartsToDto(models),
	}, nil
}
