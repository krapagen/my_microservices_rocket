package part

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, inputUuid uuid.UUID) (model.Part, error) {
	op := "Функция inventory/internl/repository/part/GetPart"
	log := slog.With("op", op)
	// 3. Найти деталь в map
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.parts[inputUuid]

	// 4. Если не найдена → NOT_FOUND

	if !ok {
		log.ErrorContext(ctx, "деталь не найдена", "uuid", inputUuid.String())
		return model.Part{}, errs.ErrPartNotFound
	}

	// 6. Вернуть деталь
	log.InfoContext(ctx, "деталь найдена", "uuid", inputUuid.String(), "name", part.Name)
	return converter.PartRecordToModel(part), nil
}
