package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func TestNewEngineProperties(t *testing.T) {
	props, err := model.NewEngineProperties(model.EngineClassC, 30)
	assert.NoError(t, err)
	assert.Equal(t, model.EngineClassC, props.Engine().Class())
	assert.Equal(t, 30, props.Engine().RequiredStrength())

	props, err = model.NewEngineProperties(model.EngineClassB, 70)
	assert.NoError(t, err)
	assert.Equal(t, 70, props.Engine().RequiredStrength())

	props, err = model.NewEngineProperties(model.EngineClassA, 100)
	assert.NoError(t, err)
	assert.Equal(t, 100, props.Engine().RequiredStrength())
}

func TestNewEngineProperties_Invalid(t *testing.T) {
	_, err := model.NewEngineProperties(model.EngineClassA, 30)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))

	_, err = model.NewEngineProperties(model.EngineClassC, 100)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))

	_, err = model.NewEngineProperties("X", 100)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))
}
