package model

import "github.com/google/uuid"

type Part struct {
	UUID          uuid.UUID
	Name          string
	PartType      PartType
	Price         int64
	StockQuantity int64
}

type PartType string

const (
	PartTypeHull   PartType = "HULL"
	PartTypeEngine PartType = "ENGINE"
	PartTypeShield PartType = "SHIELD"
	PartTypeWeapon PartType = "WEAPON"
)

// ShipSlots описывает набор UUID деталей в слотах корабля.
type ShipSlots struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}
