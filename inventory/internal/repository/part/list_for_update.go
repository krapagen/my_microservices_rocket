package part

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *repository) ListForUpdate(
	ctx context.Context,
	filter input.PartFilter,
) ([]model.Part, error) {
	return r.list(ctx, filter, "Функция inventory/internal/repository/part/ListForUpdate", true)
}
