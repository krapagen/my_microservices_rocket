package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

type PartType string

const (
	PartTypeUnspecified PartType = ""
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

func AllPartTypes() []PartType {
	return []PartType{
		PartTypeUnspecified,
		PartTypeHull,
		PartTypeEngine,
		PartTypeShield,
		PartTypeWeapon,
	}
}
