package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func TestNewShieldProperties(t *testing.T) {
	props, err := model.NewShieldProperties(model.ShieldTypeEnergy)
	assert.NoError(t, err)
	assert.Equal(t, model.ShieldTypeEnergy, props.Shield().ShieldType())

	props, err = model.NewShieldProperties(model.ShieldTypePlasma)
	assert.NoError(t, err)
	assert.Equal(t, model.ShieldTypePlasma, props.Shield().ShieldType())
}

func TestNewShieldProperties_Invalid(t *testing.T) {
	_, err := model.NewShieldProperties("unknown")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))
}

func TestShieldProperties_ConflictsWith(t *testing.T) {
	plasma, _ := model.NewShieldProperties(model.ShieldTypePlasma)
	laser, _ := model.NewWeaponProperties(model.WeaponTypeLaser)
	assert.True(t, plasma.Shield().ConflictsWith(laser.Weapon()))

	energy, _ := model.NewShieldProperties(model.ShieldTypeEnergy)
	assert.False(t, energy.Shield().ConflictsWith(laser.Weapon()))

	missile, _ := model.NewWeaponProperties(model.WeaponTypeMissile)
	assert.False(t, plasma.Shield().ConflictsWith(missile.Weapon()))
}
