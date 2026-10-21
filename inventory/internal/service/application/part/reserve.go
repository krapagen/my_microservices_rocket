package part

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Reserve увеличивает зарезервированное количество деталей на 1 для каждого указанного UUID.
// Один атомарный UPDATE в БД — без SELECT FOR UPDATE и транзакции,
// иначе конкурентные запросы сериализуются на строках parts.
func (s *service) Reserve(ctx context.Context, uuids []uuid.UUID) error {
	op := "Функция inventory/internal/service/application/part/Reserve"
	log := slog.With("op", op)
	if len(uuids) == 0 {
		log.InfoContext(ctx, "пустой список UUID")
		return nil
	}
	if err := s.partRepository.ReserveAtomic(ctx, uuids); err != nil {
		log.ErrorContext(ctx, "ошибка резервирования деталей", "error", err)
		return err
	}
	log.InfoContext(ctx, "детали зарезервированы", "count", len(uuids))
	return nil
}
