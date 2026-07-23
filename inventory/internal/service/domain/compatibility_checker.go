package domain

import (
	"fmt"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

// compatibilityChecker проверяет совместимость деталей космического корабля.
type compatibilityChecker struct{}

func NewCompatibilityChecker() *compatibilityChecker {
	return &compatibilityChecker{}
}

// Check проверяет бизнес-правила совместимости для загруженных слотов корабля.
func (c *compatibilityChecker) Check(slots model.ResolvedShipSlots) error {
	hullProperties := slots.Hull.Properties()
	engineProperties := slots.Engine.Properties()
	hullProps := hullProperties.Hull()
	engineProps := engineProperties.Engine()

	if hullProps == nil || engineProps == nil {
		return fmt.Errorf("отсутствуют свойства корпуса или двигателя: %w", errs.ErrIncompatibleParts)
	}

	if !hullProps.CanSupport(engineProps) {
		return fmt.Errorf("корпус не выдерживает двигатель: %w", errs.ErrIncompatibleParts)
	}

	if slots.Shield != nil && slots.Weapon != nil {
		shieldProperties := slots.Shield.Properties()
		weaponProperties := slots.Weapon.Properties()
		shieldProps := shieldProperties.Shield()
		weaponProps := weaponProperties.Weapon()
		if shieldProps != nil && weaponProps != nil && shieldProps.ConflictsWith(weaponProps) {
			return fmt.Errorf("щит конфликтует с оружием: %w", errs.ErrIncompatibleParts)
		}
	}

	return nil
}
