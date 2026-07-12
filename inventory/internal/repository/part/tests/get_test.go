package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func (r *RepoSuite) TestGet_Success() {
	result, err := r.repo.Get(r.ctx, uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"))

	r.NoError(err)
	r.Equal(uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), result.UUID)
	r.Equal("Алюминиевый корпус", result.Name)
	r.Equal(model.PartTypeHull, result.PartType)
	r.Equal(int64(500000), result.Price)
	r.Equal(int64(10), result.StockQuantity)
}

func (r *RepoSuite) TestGet_NotFound() {
	_, err := r.repo.Get(r.ctx, uuid.New())

	r.ErrorIs(err, errs.ErrPartNotFound)
}

func (r *RepoSuite) TestGet_AllKnownParts() {
	knownUUIDs := []struct {
		uuid     string
		name     string
		partType model.PartType
	}{
		{"550e8400-e29b-41d4-a716-446655440001", "Алюминиевый корпус", model.PartTypeHull},
		{"550e8400-e29b-41d4-a716-446655440002", "Титановый корпус", model.PartTypeHull},
		{"550e8400-e29b-41d4-a716-446655440003", "Ионный двигатель C", model.PartTypeEngine},
		{"550e8400-e29b-41d4-a716-446655440004", "Ионный двигатель B", model.PartTypeEngine},
		{"550e8400-e29b-41d4-a716-446655440005", "Энергетический щит", model.PartTypeShield},
		{"550e8400-e29b-41d4-a716-446655440006", "Лазерная пушка", model.PartTypeWeapon},
		{"550e8400-e29b-41d4-a716-446655440007", "Плазменный корпус", model.PartTypeHull},
	}

	for _, tc := range knownUUIDs {
		r.Run(tc.name, func() {
			result, err := r.repo.Get(r.ctx, uuid.MustParse(tc.uuid))

			r.NoError(err)
			r.Equal(tc.name, result.Name)
			r.Equal(tc.partType, result.PartType)
		})
	}
}

func (r *RepoSuite) TestGet_OutOfStockPart() {
	result, err := r.repo.Get(r.ctx, uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"))

	r.NoError(err)
	r.Equal("Плазменный корпус", result.Name)
	r.Equal(int64(0), result.StockQuantity)
}

func (r *RepoSuite) TestGet_EmptyUuid() {
	_, err := r.repo.Get(r.ctx, uuid.Nil)

	r.ErrorIs(err, errs.ErrPartNotFound)
}

func (r *RepoSuite) TestGet_RandomUuid() {
	_, err := r.repo.Get(r.ctx, uuid.MustParse(gofakeit.UUID()))

	r.ErrorIs(err, errs.ErrPartNotFound)
}
