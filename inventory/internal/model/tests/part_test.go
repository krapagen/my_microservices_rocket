package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func TestRestorePart(t *testing.T) {
	partUUID := uuid.New()
	createdAt := time.Now().UTC()

	part := model.RestorePart(
		partUUID,
		"Heavy Hull",
		"strong hull",
		model.PartTypeHull,
		int64(5000),
		15,
		3,
		model.PartProperties{},
		createdAt,
	)

	assert.Equal(t, partUUID, part.UUID())
	assert.Equal(t, "Heavy Hull", part.Name())
	assert.Equal(t, "strong hull", part.Description())
	assert.Equal(t, model.PartTypeHull, part.PartType())
	assert.Equal(t, int64(5000), part.Price())
	assert.Equal(t, 15, part.StockQuantity())
	assert.Equal(t, 3, part.Reserved())
	assert.Equal(t, model.PartProperties{}, part.Properties())
	assert.Equal(t, createdAt, part.CreatedAt())
}
