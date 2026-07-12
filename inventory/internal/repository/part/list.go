package part

import (
	"context"
	"log/slog"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *repository) List(
	ctx context.Context,
	filter input.PartFilter,
) ([]model.Part, error) {
	op := "Функция inventory/internl/repository/part/ListParts"
	log := slog.With("op", op)
	r.mu.RLock()
	defer r.mu.RUnlock()
	// 1. Если передан список uuids → найти детали по UUID (сохраняя порядок запроса)
	var parts []model.Part
	if len(filter.UUIDs) > 0 {
		for _, valUuid := range filter.UUIDs {

			part, ok := r.parts[valUuid]
			//    - Если хоть один UUID не найден → NOT_FOUND
			if !ok {
				log.ErrorContext(ctx, "деталь не найдена", "uuid", valUuid)
				return nil, errs.ErrPartNotFound
			}
			parts = append(parts, converter.PartRecordToModel(part))
		}
		log.InfoContext(ctx, "детали найдены по UUID", "count", len(parts))
		return parts, nil
	}

	// 2. Иначе если part_type == UNSPECIFIED → вернуть все детали
	// 3. Иначе → фильтровать по типу

	for _, part := range r.parts {
		modelPartType := converter.RecordToModelType[part.PartType]
		if filter.PartType == modelPartType || filter.PartType == model.PartTypeUnspecified {
			parts = append(parts, converter.PartRecordToModel(part))
		}
	}
	log.InfoContext(ctx, "возвращены детали по типу или все детали, если \"Unspecified\"", "count", len(parts), "part_type", filter.PartType)
	return parts, nil
}
