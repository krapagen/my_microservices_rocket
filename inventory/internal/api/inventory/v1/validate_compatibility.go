package v1

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/converter"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ValidateCompatibility(
	ctx context.Context,
	req *inventoryv1.ValidateCompatibilityRequest,
) (*inventoryv1.ValidateCompatibilityResponse, error) {
	convert := converter.NewConverter()
	slots, err := convert.ToValidateCompatibilityInput(req)
	if err != nil {
		return nil, err
	}

	if err := a.partService.ValidateCompatibility(ctx, slots); err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
