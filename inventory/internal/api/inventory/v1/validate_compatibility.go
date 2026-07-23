package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/converter"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ValidateCompatibility(
	ctx context.Context,
	req *inventoryv1.ValidateCompatibilityRequest,
) (*inventoryv1.ValidateCompatibilityResponse, error) {
	convert := converter.NewConverter()

	toUUID := func(raw string) uuid.UUID {
		if raw == "" {
			return uuid.Nil
		}
		u, err := convert.ToGetInput(raw)
		if err != nil {
			return uuid.Nil
		}
		return u
	}

	slots := model.ShipSlots{
		HullUUID:   toUUID(req.GetHullUuid()),
		EngineUUID: toUUID(req.GetEngineUuid()),
		ShieldUUID: toUUID(req.GetShieldUuid()),
		WeaponUUID: toUUID(req.GetWeaponUuid()),
	}

	if err := a.partService.ValidateCompatibility(ctx, slots); err != nil {
		return nil, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
