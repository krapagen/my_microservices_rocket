package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func TestNewWeaponProperties(t *testing.T) {
	props, err := model.NewWeaponProperties(model.WeaponTypeLaser)
	assert.NoError(t, err)
	assert.Equal(t, model.WeaponTypeLaser, props.Weapon().WeaponType())

	props, err = model.NewWeaponProperties(model.WeaponTypeMissile)
	assert.NoError(t, err)
	assert.Equal(t, model.WeaponTypeMissile, props.Weapon().WeaponType())
}

func TestNewWeaponProperties_Invalid(t *testing.T) {
	_, err := model.NewWeaponProperties("unknown")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))
}
