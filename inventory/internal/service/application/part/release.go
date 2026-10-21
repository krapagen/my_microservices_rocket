package part

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Release уменьшает зарезервированное количество деталей на 1 для каждого указанного UUID.
// Один атомарный UPDATE — без SELECT FOR UPDATE и транзакции.
func (s *service) Release(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/service/application/part/Release"
	log := slog.With("op", op)
	if len(uuids) == 0 {
		log.InfoContext(ctx, "пустой список UUID")
		return nil
	}
	if err := s.partRepository.ReleaseAtomic(ctx, uuids); err != nil {
		log.ErrorContext(ctx, "ошибка освобождения деталей", "error", err)
		return err
	}
	log.InfoContext(ctx, "детали освобождены", "count", len(uuids))
	return nil
}
