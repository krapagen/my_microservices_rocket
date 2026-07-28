package converter

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (c *converter) ToValidateCompatibilityInput(req *inventoryv1.ValidateCompatibilityRequest) (model.ShipSlots, error) {
	if req == nil {
		return model.ShipSlots{}, nil
	}

	hullUUID, err := c.toOptionalUUID(req.GetHullUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	engineUUID, err := c.toOptionalUUID(req.GetEngineUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	shieldUUID, err := c.toOptionalUUID(req.GetShieldUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	weaponUUID, err := c.toOptionalUUID(req.GetWeaponUuid())
	if err != nil {
		return model.ShipSlots{}, err
	}

	return model.ShipSlots{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}, nil
}

func (c *converter) toOptionalUUID(raw string) (uuid.UUID, error) {
	if raw == "" {
		return uuid.Nil, nil
	}

	return c.ToGetInput(raw)
}
