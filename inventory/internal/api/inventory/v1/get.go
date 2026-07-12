package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/converter"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	op := "Функция inventory/internal/api/inventory/v1/GetPart"
	log := slog.With("op", op)

	convert := converter.NewConverter()
	reqUuid, err := convert.ToGetInput(req.GetUuid())
	if err != nil {
		log.ErrorContext(ctx, "неверный формат uuid", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "валидный формат uuid", "uuid", reqUuid)
	modelPart, err := a.partService.Get(ctx, reqUuid)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при получении детали", "error", err)
		return nil, err
	}

	log.InfoContext(ctx, "деталь найдена", "uuid", modelPart.UUID.String(), "name", modelPart.Name)

	return &inventoryv1.GetPartResponse{
		Part: convert.PartToDto(modelPart),
	}, nil
}
