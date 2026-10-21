package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

// ReleaseAtomic декрементирует reserved по всем UUID в одной транзакции.
// По-строчные UPDATE в отсортированном порядке (как ReserveAtomic/Commit) —
// канонический порядок локов исключает дедлоки, откат tx сохраняет
// all-or-nothing.
func (r *repository) ReleaseAtomic(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/repository/part/ReleaseAtomic"
	log := slog.With("op", op)

	sorted := sortedUUIDs(uuids)
	return r.txMgr.Do(ctx, func(txCtx context.Context) error {
		// CTE + ORDER BY FOR UPDATE: канонический порядок локов, один
		// round-trip — см. ReserveAtomic.
		const query = `
			WITH locked AS (
				SELECT uuid FROM parts
				WHERE uuid = ANY($1::uuid[])
				ORDER BY uuid
				FOR UPDATE
			)
			UPDATE parts p
			SET reserved = p.reserved - 1
			FROM locked l
			WHERE p.uuid = l.uuid
			  AND p.reserved > 0`
		res, err := r.getter.DefaultTrOrDB(txCtx, r.pool).Exec(txCtx, query, sorted)
		if err != nil {
			log.ErrorContext(txCtx, "ошибка освобождения деталей", "error", err)
			return fmt.Errorf("release parts: %w", err)
		}
		if res.RowsAffected() != int64(len(sorted)) {
			return r.releaseFailureReason(txCtx, sorted)
		}
		return nil
	})
}

// releaseFailureReason вызывается только на пути ошибки: возвращает
// ErrPartNotFound (детали нет) или ErrNothingToRelease (reserved == 0).
func (r *repository) releaseFailureReason(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/repository/part/releaseFailureReason"
	log := slog.With("op", op)

	const query = `
		SELECT uuid, reserved
		FROM parts
		WHERE uuid = ANY($1::uuid[])
		FOR UPDATE`
	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, uuids)
	if err != nil {
		log.ErrorContext(ctx, "ошибка чтения деталей для диагностики освобождения", "error", err)
		return fmt.Errorf("list parts for release diagnostics: %w", err)
	}
	type releaseRow struct {
		UUID     uuid.UUID `db:"uuid"`
		Reserved int       `db:"reserved"`
	}
	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[releaseRow])
	if err != nil {
		log.ErrorContext(ctx, "ошибка конвертации деталей для диагностики", "error", err)
		return fmt.Errorf("collect parts for release diagnostics: %w", err)
	}

	found := make(map[uuid.UUID]struct{}, len(records))
	for _, rec := range records {
		found[rec.UUID] = struct{}{}
		if rec.Reserved <= 0 {
			return fmt.Errorf("деталь %s не имеет зарезервированных единиц: %w", rec.UUID, errs.ErrNothingToRelease)
		}
	}
	for _, id := range uuids {
		if _, ok := found[id]; !ok {
			return fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, id)
		}
	}
	return errs.ErrNothingToRelease
}
