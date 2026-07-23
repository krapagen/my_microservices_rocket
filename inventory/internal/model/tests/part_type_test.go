package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func TestNewPartType(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected model.PartType
	}{
		{"HULL", model.PartTypeHull},
		{"ENGINE", model.PartTypeEngine},
		{"SHIELD", model.PartTypeShield},
		{"WEAPON", model.PartTypeWeapon},
	} {
		got, err := model.NewPartType(tt.input)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, got)
	}

	_, err := model.NewPartType("UNKNOWN")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProperties))
}

func TestAllPartTypes(t *testing.T) {
	assert.Equal(t,
		[]model.PartType{
			model.PartTypeUnspecified,
			model.PartTypeHull,
			model.PartTypeEngine,
			model.PartTypeShield,
			model.PartTypeWeapon,
		},
		model.AllPartTypes(),
	)
}
