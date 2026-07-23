package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

type PartService interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
	Reserve(ctx context.Context, uuids []uuid.UUID) error
	Release(ctx context.Context, uuids []uuid.UUID) error
	ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error
}
