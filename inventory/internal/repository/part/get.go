package part

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, inputUuid uuid.UUID) (model.Part, error) {
	op := "Функция inventory/internl/repository/part/GetPart"
	log := slog.With("op", op)
	const query = `
		SELECT 
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.created_at
		FROM parts AS p
		WHERE p.uuid = $1;`

	var part record.Part

	err := r.pool.QueryRow(ctx, query, inputUuid).Scan(
		&part.UUID,
		&part.Name,
		&part.Description,
		&part.PartType,
		&part.Price,
		&part.StockQuantity,
		&part.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.ErrorContext(ctx, "деталь не найдена", "error", err)
			return model.Part{}, fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, inputUuid)
		}
		log.InfoContext(ctx, "ошибка получения детали", "error", err)
		return model.Part{}, err
	}
	log.InfoContext(ctx, "деталь успешно получена", "uuid", inputUuid)
	return converter.PartRecordToModel(part), nil
}
