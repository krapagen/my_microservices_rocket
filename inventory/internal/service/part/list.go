package part

import (
	"context"
	"log/slog"
	"sort"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *service) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	op := "Функция inventory/internl/service/part/List"
	log := slog.With("op", op)
	parts, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Сортируем по имени
	if len(filter.UUIDs) == 0 {
		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Name < parts[j].Name
		})
		log.InfoContext(ctx, "возвращены отсортированные детали", "count", len(parts))
	}

	return parts, nil
}
