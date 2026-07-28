package part

import (
	"context"
	"errors"
	"log/slog"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *service) Commit(ctx context.Context, filter input.CommitFilter) error {
	op := "Функция inventory/internal/service/application/part/Commit"
	log := slog.With("op", op)

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		parts, err := s.partRepository.ListForUpdate(ctx, input.PartFilter{UUIDs: filter.UUIDs})
		if err != nil {
			log.ErrorContext(ctx, "ошибка получения деталей для обновления", "error", err)
			return err
		}
		log.InfoContext(ctx, "успешно прочитаны детали для обновления", "count", len(parts))
		for _, p := range parts {
			if p.StockQuantity() <= 0 {
				log.ErrorContext(ctx, "деталь отсутствует на складе:", "uuid", p.UUID())
				return errs.ErrCommitParts
			}
			if p.Reserved() <= 0 {
				log.ErrorContext(ctx, "деталь не зарезервирована:", "uuid", p.UUID())
				return errs.ErrCommitParts
			}
		}
		err = s.partRepository.Commit(ctx, filter)
		if err != nil {
			if errors.Is(err, errs.ErrCommitParts) {
				log.ErrorContext(ctx, "не удалось обновить все детали", "expected", len(filter.UUIDs), "actual", len(parts))
				return errs.ErrCommitParts
			}
			log.ErrorContext(ctx, "ошибка при списании деталей", "error", err)
			return err
		}
		log.InfoContext(ctx, "успешно списаны детали", "count", len(parts))
		return nil
	})
	if err != nil {
		log.ErrorContext(ctx, "ошибка при выполнении транзакции", "error", err)
		return err
	}
	return nil
}
