package part

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *service) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	parts, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Сортируем по имени

	return parts, nil
}
