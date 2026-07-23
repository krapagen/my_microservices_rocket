package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	return s.partRepository.Get(ctx, partUUID)
}
