package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *repository) Commit(
	ctx context.Context,
	filter input.CommitFilter,
) error {
	op := "Функция inventory/internal/repository/part/Commit"
	log := slog.With("op", op)
	// CTE + ORDER BY FOR UPDATE: канонический порядок локов, один
	// round-trip — см. ReserveAtomic. Откат tx при частичном обновлении
	// сохраняет all-or-nothing.
	sorted := sortedUUIDs(filter.UUIDs)
	return r.txMgr.Do(ctx, func(txCtx context.Context) error {
		query := `WITH locked AS (
				SELECT uuid FROM parts
				WHERE uuid = ANY($1::uuid[])
				ORDER BY uuid
				FOR UPDATE
			)
			UPDATE parts p
			SET stock_quantity = p.stock_quantity - 1,
				reserved = p.reserved - 1
			FROM locked l
			WHERE p.uuid = l.uuid
			  AND p.stock_quantity > 0
			  AND p.reserved > 0`
		res, err := r.getter.DefaultTrOrDB(txCtx, r.pool).Exec(txCtx, query, sorted)
		if err != nil {
			log.ErrorContext(txCtx, "Ошибка при выполнении запроса", "error", err)
			return fmt.Errorf("commit parts: %w", err)
		}
		if res.RowsAffected() != int64(len(sorted)) {
			log.ErrorContext(txCtx, "Не удалось обновить все детали", "want", len(sorted), "got", res.RowsAffected())
			return r.commitFailureReason(txCtx, sorted)
		}
		return nil
	})
}

// commitFailureReason вызывается только на пути ошибки: возвращает
// ErrPartNotFound (детали нет) или ErrCommitParts (есть, но
// stock_quantity/reserved == 0 — списывать нечего).
func (r *repository) commitFailureReason(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/repository/part/commitFailureReason"
	log := slog.With("op", op)

	const query = `
		SELECT uuid
		FROM parts
		WHERE uuid = ANY($1::uuid[])
		FOR UPDATE`
	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, uuids)
	if err != nil {
		log.ErrorContext(ctx, "ошибка чтения деталей для диагностики списания", "error", err)
		return fmt.Errorf("list parts for commit diagnostics: %w", err)
	}
	type commitRow struct {
		UUID uuid.UUID `db:"uuid"`
	}
	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[commitRow])
	if err != nil {
		log.ErrorContext(ctx, "ошибка конвертации деталей для диагностики", "error", err)
		return fmt.Errorf("collect parts for commit diagnostics: %w", err)
	}

	found := make(map[uuid.UUID]struct{}, len(records))
	for _, rec := range records {
		found[rec.UUID] = struct{}{}
	}
	for _, id := range uuids {
		if _, ok := found[id]; !ok {
			return fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, id)
		}
	}
	return fmt.Errorf("обновить строки при списании не удалось: %w", errs.ErrCommitParts)
}
