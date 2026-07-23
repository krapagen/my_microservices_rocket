package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func TestNewHullProperties(t *testing.T) {
	props, err := model.NewHullProperties(100)
	assert.NoError(t, err)
	assert.Equal(t, 100, props.Hull().Strength())
}

func TestNewHullProperties_Invalid(t *testing.T) {
	_, err := model.NewHullProperties(10)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))

	_, err = model.NewHullProperties(250)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))
}

func TestHullProperties_CanSupport(t *testing.T) {
	hull, _ := model.NewHullProperties(100)
	engine, _ := model.NewEngineProperties(model.EngineClassA, 100)
	assert.True(t, hull.Hull().CanSupport(engine.Engine()))

	weakHull, _ := model.NewHullProperties(50)
	assert.False(t, weakHull.Hull().CanSupport(engine.Engine()))

	engineC, _ := model.NewEngineProperties(model.EngineClassC, 30)
	assert.True(t, weakHull.Hull().CanSupport(engineC.Engine()))
}
