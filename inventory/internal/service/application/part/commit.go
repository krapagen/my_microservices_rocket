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

	// Репозиторный UPDATE атомарен: WHERE stock>0 AND reserved>0 + сверка
	// RowsAffected заменяют SELECT FOR UPDATE pre-check — лок на время UPDATE.
	if err := s.partRepository.Commit(ctx, filter); err != nil {
		if errors.Is(err, errs.ErrCommitParts) {
			log.ErrorContext(ctx, "не удалось списать все детали", "expected", len(filter.UUIDs))
			return errs.ErrCommitParts
		}
		log.ErrorContext(ctx, "ошибка при списании деталей", "error", err)
		return err
	}
	log.InfoContext(ctx, "успешно списаны детали", "count", len(filter.UUIDs))
	return nil
}
