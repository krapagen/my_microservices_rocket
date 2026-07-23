package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/domain"
)

func newPart(t model.PartType, props model.PartProperties) model.Part {
	return model.RestorePart(uuid.New(), "part", "", t, 1000, 10, 0, props, time.Now().UTC())
}

func TestCompatibilityChecker_Success(t *testing.T) {
	hullProps, _ := model.NewHullProperties(100)
	engineProps, _ := model.NewEngineProperties(model.EngineClassA, 100)
	slots := model.ResolvedShipSlots{
		Hull:   newPart(model.PartTypeHull, hullProps),
		Engine: newPart(model.PartTypeEngine, engineProps),
	}

	checker := domain.NewCompatibilityChecker()
	err := checker.Check(slots)
	assert.NoError(t, err)
}

func TestCompatibilityChecker_SuccessWithOptional(t *testing.T) {
	hullProps, _ := model.NewHullProperties(100)
	engineProps, _ := model.NewEngineProperties(model.EngineClassC, 30)
	shieldProps, _ := model.NewShieldProperties(model.ShieldTypeEnergy)
	weaponProps, _ := model.NewWeaponProperties(model.WeaponTypeMissile)
	slots := model.ResolvedShipSlots{
		Hull:   newPart(model.PartTypeHull, hullProps),
		Engine: newPart(model.PartTypeEngine, engineProps),
		Shield: &[]model.Part{newPart(model.PartTypeShield, shieldProps)}[0],
		Weapon: &[]model.Part{newPart(model.PartTypeWeapon, weaponProps)}[0],
	}

	checker := domain.NewCompatibilityChecker()
	err := checker.Check(slots)
	assert.NoError(t, err)
}

func TestCompatibilityChecker_HullTooWeak(t *testing.T) {
	hullProps, _ := model.NewHullProperties(50)
	engineProps, _ := model.NewEngineProperties(model.EngineClassA, 100)
	slots := model.ResolvedShipSlots{
		Hull:   newPart(model.PartTypeHull, hullProps),
		Engine: newPart(model.PartTypeEngine, engineProps),
	}

	checker := domain.NewCompatibilityChecker()
	err := checker.Check(slots)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrIncompatibleParts))
}

func TestCompatibilityChecker_PlasmaLaserConflict(t *testing.T) {
	hullProps, _ := model.NewHullProperties(100)
	engineProps, _ := model.NewEngineProperties(model.EngineClassC, 30)
	shieldProps, _ := model.NewShieldProperties(model.ShieldTypePlasma)
	weaponProps, _ := model.NewWeaponProperties(model.WeaponTypeLaser)
	slots := model.ResolvedShipSlots{
		Hull:   newPart(model.PartTypeHull, hullProps),
		Engine: newPart(model.PartTypeEngine, engineProps),
		Shield: &[]model.Part{newPart(model.PartTypeShield, shieldProps)}[0],
		Weapon: &[]model.Part{newPart(model.PartTypeWeapon, weaponProps)}[0],
	}

	checker := domain.NewCompatibilityChecker()
	err := checker.Check(slots)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrIncompatibleParts))
}

func TestCompatibilityChecker_MissingHullProperties(t *testing.T) {
	engineProps, _ := model.NewEngineProperties(model.EngineClassC, 30)
	slots := model.ResolvedShipSlots{
		Hull:   newPart(model.PartTypeHull, model.PartProperties{}),
		Engine: newPart(model.PartTypeEngine, engineProps),
	}

	checker := domain.NewCompatibilityChecker()
	err := checker.Check(slots)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrIncompatibleParts))
}

func TestCompatibilityChecker_MissingEngineProperties(t *testing.T) {
	hullProps, _ := model.NewHullProperties(100)
	slots := model.ResolvedShipSlots{
		Hull:   newPart(model.PartTypeHull, hullProps),
		Engine: newPart(model.PartTypeEngine, model.PartProperties{}),
	}

	checker := domain.NewCompatibilityChecker()
	err := checker.Check(slots)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrIncompatibleParts))
}
