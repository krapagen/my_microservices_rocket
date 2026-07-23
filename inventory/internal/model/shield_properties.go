package model

import (
	"fmt"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

// ShieldType — тип щита.
type ShieldType string

const (
	ShieldTypeEnergy ShieldType = "energy"
	ShieldTypePlasma ShieldType = "plasma"
)

// ShieldProperties — свойства щита (Value Object).
type ShieldProperties struct {
	shieldType ShieldType
}

func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

// ConflictsWith проверяет, создаёт ли щит помехи оружию.
// Плазменный щит конфликтует с лазерным оружием.
func (s *ShieldProperties) ConflictsWith(w *WeaponProperties) bool {
	return s.shieldType == ShieldTypePlasma && w.WeaponType() == WeaponTypeLaser
}

// NewShieldProperties создаёт свойства щита.
func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	switch shieldType {
	case ShieldTypeEnergy, ShieldTypePlasma:
		return PartProperties{
			shield: &ShieldProperties{shieldType: shieldType},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип щита %q: %w", shieldType, errs.ErrInvalidProperties)
	}
}
