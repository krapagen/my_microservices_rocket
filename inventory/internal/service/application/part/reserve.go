package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

// Reserve увеличивает зарезервированное количество деталей на 1 для каждого указанного UUID.
func (s *service) Reserve(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/service/application/part/Reserve"
	log := slog.With("op", op)
	if len(uuids) == 0 {
		log.InfoContext(ctx, "пустой список UUID")
		return nil
	}

	parts, err := s.partRepository.List(ctx, input.PartFilter{UUIDs: uuids})
	if err != nil {
		log.ErrorContext(ctx, "ошибка получения деталей", "error", err)
		return err
	}
	log.InfoContext(ctx, "успешно прочитаны детали", "count", len(parts))
	updated := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		if p.Reserved() >= p.StockQuantity() {
			log.ErrorContext(ctx, "деталь %s отсутствует на складе", "uuid", p.UUID())
			return fmt.Errorf("деталь %s отсутствует на складе: %w", p.UUID(), errs.ErrOutOfStock)
		}

		updated = append(updated, model.RestorePart(
			p.UUID(),
			p.Name(),
			p.Description(),
			p.PartType(),
			p.Price(),
			p.StockQuantity(),
			p.Reserved()+1,
			p.Properties(),
			p.CreatedAt(),
		))
	}
	log.InfoContext(ctx, "успешно подготовлены детали для обновления зарезервированного количества", "count", len(updated))
	return s.partRepository.UpdateReservedBatch(ctx, updated)
}
