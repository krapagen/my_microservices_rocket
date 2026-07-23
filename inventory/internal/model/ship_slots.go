package model

import "github.com/google/uuid"

type ShipSlots struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID uuid.UUID
	WeaponUUID uuid.UUID
}

type ResolvedShipSlots struct {
	Hull   Part
	Engine Part
	Shield *Part
	Weapon *Part
}
