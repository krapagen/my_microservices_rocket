package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) CommitParts(
	ctx context.Context,
	req *inventoryv1.CommitPartsRequest,
) (*inventoryv1.CommitPartsResponse, error) {
	op := "Функция inventory/internal/api/inventory/v1/CommitParts"
	log := slog.With("op", op)
	convert := converter.NewConverter()
	commitFilter, err := convert.ToGetInputs(req.GetUuids())
	if err != nil {
		log.ErrorContext(ctx, "неверный формат uuid", "error", err)
		return nil, err
	}

	err = a.partService.Commit(ctx, input.CommitFilter{UUIDs: commitFilter})
	if err != nil {
		log.ErrorContext(ctx, "ошибка при коммите частей", "error", err)
		return nil, err
	}
	return &inventoryv1.CommitPartsResponse{}, nil
}
