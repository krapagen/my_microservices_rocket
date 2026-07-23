package model

import (
	"fmt"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

// WeaponType — тип оружия.
type WeaponType string

const (
	WeaponTypeLaser   WeaponType = "laser"
	WeaponTypeMissile WeaponType = "missile"
)

// WeaponProperties — свойства оружия (Value Object).
type WeaponProperties struct {
	weaponType WeaponType
}

func (w *WeaponProperties) WeaponType() WeaponType { return w.weaponType }

// NewWeaponProperties создаёт свойства оружия.
func NewWeaponProperties(weaponType WeaponType) (PartProperties, error) {
	switch weaponType {
	case WeaponTypeLaser, WeaponTypeMissile:
		return PartProperties{
			weapon: &WeaponProperties{weaponType: weaponType},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("неизвестный тип оружия %q: %w", weaponType, errs.ErrInvalidProperties)
	}
}
