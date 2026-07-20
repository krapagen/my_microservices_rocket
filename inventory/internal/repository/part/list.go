package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/record"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *repository) List(
	ctx context.Context,
	filter input.PartFilter,
) ([]model.Part, error) {
	query := `
		SELECT
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.created_at
		FROM parts AS p`

	var args []any

	switch {
	case len(filter.UUIDs) > 0:
		query += ` WHERE p.uuid = ANY($1::uuid[])`
		args = append(args, filter.UUIDs)
	case filter.PartType != model.PartTypeUnspecified:
		query += ` WHERE p.part_type = $1
                   ORDER BY p.name ASC`
		args = append(args, string(filter.PartType))
	default:
		query += ` ORDER BY p.name ASC`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list parts: %w", err)
	}
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, fmt.Errorf("collect parts: %w", err)
	}

	if len(filter.UUIDs) > 0 {
		partsByUUID := make(map[uuid.UUID]model.Part, len(records))
		for _, rec := range records {
			partsByUUID[rec.UUID] = converter.PartRecordToModel(rec)
		}

		result := make([]model.Part, 0, len(filter.UUIDs))
		for _, id := range filter.UUIDs {
			part, ok := partsByUUID[id]
			if !ok {
				return nil, fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, id)
			}
			result = append(result, part)
		}
		return result, nil
	}

	return converter.PartsRecordToModel(records), nil
}
