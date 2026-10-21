package part

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

// ReserveAtomic инкрементирует reserved по всем UUID в одной транзакции.
// UPDATE идёт по одному UUID за раз и всегда в отсортированном порядке:
// канонический порядок захвата row-lock'ов математически исключает
// дедлоки между конкурентными батчами ({A,B} vs {B,A}), а all-or-nothing
// обеспечивает откат транзакции при первой же ошибке.
//
// Условие reserved < stock_quantity внутри запроса делает проверку
// остатка атомарной — без предварительного SELECT FOR UPDATE.
//
// Если какая-то деталь не обновилась — возвращает ErrPartNotFound
// (детали нет) или ErrOutOfStock (деталь есть, но reserved >= stock).
func (r *repository) ReserveAtomic(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/repository/part/ReserveAtomic"
	log := slog.With("op", op)

	sorted := sortedUUIDs(uuids)
	return r.txMgr.Do(ctx, func(txCtx context.Context) error {
		// CTE захватывает row-локи в отсортированном порядке (канонический
		// порядок — дедлоки исключены), внешний UPDATE идёт по уже залоченным
		// строкам без новых ожиданий. Один statement = один round-trip:
		// лок удерживается минимальное время — критично для горячих строк,
		// где стоимость tx определяет ёмкость очереди ожидания.
		const query = `
			WITH locked AS (
				SELECT uuid FROM parts
				WHERE uuid = ANY($1::uuid[])
				ORDER BY uuid
				FOR UPDATE
			)
			UPDATE parts p
			SET reserved = p.reserved + 1
			FROM locked l
			WHERE p.uuid = l.uuid
			  AND p.reserved < p.stock_quantity`
		res, err := r.getter.DefaultTrOrDB(txCtx, r.pool).Exec(txCtx, query, sorted)
		if err != nil {
			log.ErrorContext(txCtx, "ошибка резервирования деталей", "error", err)
			return fmt.Errorf("reserve parts: %w", err)
		}
		if res.RowsAffected() != int64(len(sorted)) {
			return r.reserveFailureReason(txCtx, sorted)
		}
		return nil
	})
}

// sortedUUIDs возвращает отсортированную копию массива UUID.
// Единый канонический порядок обхода строк для всех батч-операций —
// условие отсутствия дедлоков между конкурентными UPDATE.
func sortedUUIDs(uuids []uuid.UUID) []uuid.UUID {
	sorted := make([]uuid.UUID, len(uuids))
	copy(sorted, uuids)
	slices.SortFunc(sorted, func(a, b uuid.UUID) int {
		return strings.Compare(a.String(), b.String())
	})
	return sorted
}

// reserveFailureReason вызывается только на пути ошибки: дочитывает детали
// под FOR UPDATE и возвращает предметную ошибку.
func (r *repository) reserveFailureReason(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/repository/part/reserveFailureReason"
	log := slog.With("op", op)

	const query = `
		SELECT uuid, stock_quantity, reserved
		FROM parts
		WHERE uuid = ANY($1::uuid[])
		FOR UPDATE`
	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, uuids)
	if err != nil {
		log.ErrorContext(ctx, "ошибка чтения деталей для диагностики резервирования", "error", err)
		return fmt.Errorf("list parts for reserve diagnostics: %w", err)
	}
	type stockRow struct {
		UUID          uuid.UUID `db:"uuid"`
		StockQuantity int       `db:"stock_quantity"`
		Reserved      int       `db:"reserved"`
	}
	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[stockRow])
	if err != nil {
		log.ErrorContext(ctx, "ошибка конвертации деталей для диагностики", "error", err)
		return fmt.Errorf("collect parts for reserve diagnostics: %w", err)
	}

	found := make(map[uuid.UUID]struct{}, len(records))
	for _, rec := range records {
		found[rec.UUID] = struct{}{}
		if rec.Reserved >= rec.StockQuantity {
			return fmt.Errorf("деталь %s отсутствует на складе: %w", rec.UUID, errs.ErrOutOfStock)
		}
	}
	for _, id := range uuids {
		if _, ok := found[id]; !ok {
			return fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, id)
		}
	}
	return errs.ErrOutOfStock
}
