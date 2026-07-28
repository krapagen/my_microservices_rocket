package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/record"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *repository) list(
	ctx context.Context,
	filter input.PartFilter,
	op string,
	forUpdate bool,
) ([]model.Part, error) {
	log := slog.With("op", op)

	query := `
		SELECT
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.reserved,
			p.properties,
			p.created_at
		FROM parts AS p`

	var args []any

	suffix := ""
	if forUpdate {
		suffix = ` FOR UPDATE`
	}

	switch {
	case len(filter.UUIDs) > 0:
		query += ` WHERE p.uuid = ANY($1::uuid[])` + suffix
		args = append(args, filter.UUIDs)
	case filter.PartType != model.PartTypeUnspecified:
		query += ` WHERE p.part_type = $1
                   ORDER BY p.name ASC` + suffix
		args = append(args, string(filter.PartType))
	default:
		query += ` ORDER BY p.name ASC` + suffix
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, args...)
	if err != nil {
		log.ErrorContext(ctx, "ошибка получения списка деталей", "error", err)
		return nil, fmt.Errorf("list parts: %w", err)
	}
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.PartRecord])
	if err != nil {
		log.ErrorContext(ctx, "ошибка конвертации списка деталей", "error", err)
		return nil, fmt.Errorf("collect parts: %w", err)
	}
	log.InfoContext(ctx, "список деталей успешно конвертирован в record.PartRecord")

	if len(filter.UUIDs) > 0 {
		partsByUUID := make(map[uuid.UUID]model.Part, len(records))
		for _, rec := range records {
			part, err := converter.PartRecordToModel(rec)
			if err != nil {
				log.ErrorContext(ctx, "ошибка конвертации детали", "error", err)
				return nil, fmt.Errorf("convert part uuid=%s: %w", rec.UUID, err)
			}
			partsByUUID[rec.UUID] = part
		}
		log.InfoContext(ctx, "список деталей успешно конвертирован в model.Part")
		result := make([]model.Part, 0, len(filter.UUIDs))
		for _, id := range filter.UUIDs {
			part, ok := partsByUUID[id]
			if !ok {
				log.ErrorContext(ctx, "деталь не найдена", "uuid", id)
				return nil, fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, id)
			}
			result = append(result, part)
		}
		log.InfoContext(ctx, "список деталей успешно сформирован по UUID")
		return result, nil
	}

	result, err := converter.PartsRecordToModel(records)
	if err != nil {
		log.ErrorContext(ctx, "ошибка конвертации списка деталей", "error", err)
		return nil, fmt.Errorf("convert parts: %w", err)
	}
	log.InfoContext(ctx, "список деталей успешно конвертирован")

	return result, nil
}
