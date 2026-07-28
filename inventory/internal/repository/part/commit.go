package part

import (
	"context"
	"fmt"
	"log/slog"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *repository) Commit(
	ctx context.Context,
	filter input.CommitFilter,
) error {
	op := "Функция inventory/internal/repository/part/Commit"
	log := slog.With("op", op)
	query := `UPDATE parts
		SET stock_quantity = stock_quantity - 1,
			reserved = reserved - 1
		WHERE uuid = ANY($1)
		  AND stock_quantity > 0
		  AND reserved > 0`
	res, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, filter.UUIDs)
	if err != nil {
		log.ErrorContext(ctx, "Ошибка при выполнении запроса", "error", err)
		return fmt.Errorf("commit parts: %w", err)
	}
	if res.RowsAffected() < int64(len(filter.UUIDs)) {
		log.ErrorContext(ctx, "Не удалось обновить все детали", "expected", len(filter.UUIDs), "actual", res.RowsAffected())
		return fmt.Errorf("обновить строки при списании не удалось: %w", errs.ErrCommitParts)
	}
	return nil
}
