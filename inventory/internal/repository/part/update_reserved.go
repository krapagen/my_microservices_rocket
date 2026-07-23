package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

// UpdateReservedBatch обновляет поле reserved для нескольких деталей за один round-trip к БД.
func (r *repository) UpdateReservedBatch(ctx context.Context, parts []model.Part) error {
	op := "Функция inventory/internl/repository/part/UpdateReservedBatch"
	log := slog.With("op", op)
	if len(parts) == 0 {
		log.WarnContext(ctx, "пустой список деталей для обновления reserved")
		return nil
	}

	ids := make([]uuid.UUID, len(parts))
	reserved := make([]int, len(parts))
	for i, p := range parts {
		ids[i] = p.UUID()
		reserved[i] = p.Reserved()
	}

	const query = `
		UPDATE parts AS p
		SET reserved = batch.reserved
		FROM unnest($1::uuid[], $2::int[]) AS batch(uuid, reserved)
		WHERE p.uuid = batch.uuid`

	res, err := r.pool.Exec(ctx, query, ids, reserved)
	if err != nil {
		log.ErrorContext(ctx, "ошибка обновления reserved", "error", err)
		return fmt.Errorf("update reserved: %w", err)
	}
	if res.RowsAffected() == 0 {
		log.WarnContext(ctx, "не удалось обновить reserved ни для одной детали")
	}
	log.InfoContext(ctx, "reserved успешно обновлен")
	return nil
}
