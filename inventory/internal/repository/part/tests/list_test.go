package tests

import (
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (r *RepoSuite) TestList_AllParts() {
	result, err := r.repo.List(r.ctx, input.PartFilter{PartType: model.PartTypeUnspecified})

	r.NoError(err)
	r.Len(result, 7)
}

func (r *RepoSuite) TestList_ZeroParts() {
	result, err := r.repo.List(r.ctx, input.PartFilter{})

	r.NoError(err)
	r.Len(result, 7)
}

func (r *RepoSuite) TestList_ByType_Hull() {
	result, err := r.repo.List(r.ctx, input.PartFilter{PartType: model.PartTypeHull})

	r.NoError(err)
	r.Len(result, 3)
	for _, p := range result {
		r.Equal(model.PartTypeHull, p.PartType)
	}
}

func (r *RepoSuite) TestList_ByType_Engine() {
	result, err := r.repo.List(r.ctx, input.PartFilter{PartType: model.PartTypeEngine})

	r.NoError(err)
	r.Len(result, 2)
	for _, p := range result {
		r.Equal(model.PartTypeEngine, p.PartType)
	}
}

func (r *RepoSuite) TestList_ByType_Shield() {
	result, err := r.repo.List(r.ctx, input.PartFilter{PartType: model.PartTypeShield})

	r.NoError(err)
	r.Len(result, 1)
	r.Equal(model.PartTypeShield, result[0].PartType)
}

func (r *RepoSuite) TestList_ByType_Weapon() {
	result, err := r.repo.List(r.ctx, input.PartFilter{PartType: model.PartTypeWeapon})

	r.NoError(err)
	r.Len(result, 1)
	r.Equal(model.PartTypeWeapon, result[0].PartType)
}

func (r *RepoSuite) TestList_ByUUIDs_Success() {
	uuids := []uuid.UUID{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
	}

	result, err := r.repo.List(r.ctx, input.PartFilter{UUIDs: uuids})

	r.NoError(err)
	r.Len(result, 2)
	r.Equal(uuids[0], result[0].UUID)
	r.Equal(uuids[1], result[1].UUID)
}

func (r *RepoSuite) TestList_ByUUIDs_PreservesOrder() {
	uuids := []uuid.UUID{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
	}

	result, err := r.repo.List(r.ctx, input.PartFilter{UUIDs: uuids})

	r.NoError(err)
	r.Len(result, 3)
	r.Equal(uuids[0], result[0].UUID)
	r.Equal(uuids[1], result[1].UUID)
	r.Equal(uuids[2], result[2].UUID)
}

func (r *RepoSuite) TestList_ByUUIDs_NotFound() {
	uuids := []uuid.UUID{uuid.New()}

	_, err := r.repo.List(r.ctx, input.PartFilter{UUIDs: uuids})

	r.ErrorIs(err, errs.ErrPartNotFound)
}

func (r *RepoSuite) TestList_ByUUIDs_1NotFound() {
	uuids := []uuid.UUID{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
		uuid.New(),
	}

	_, err := r.repo.List(r.ctx, input.PartFilter{UUIDs: uuids})

	r.ErrorIs(err, errs.ErrPartNotFound)
}

func (r *RepoSuite) TestList_ByUUIDs_PreservesOrder_WithPartFilter() {
	uuids := []uuid.UUID{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
	}

	result, err := r.repo.List(r.ctx, input.PartFilter{UUIDs: uuids, PartType: model.PartTypeShield})

	r.NoError(err)
	r.Len(result, 3)
	r.Equal(uuids[0], result[0].UUID)
	r.Equal(uuids[1], result[1].UUID)
	r.Equal(uuids[2], result[2].UUID)
}

func (r *RepoSuite) TestList_ByUUIDs_PreservesOrder_WithPartFilterUnspecified() {
	uuids := []uuid.UUID{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
	}

	result, err := r.repo.List(r.ctx, input.PartFilter{UUIDs: uuids, PartType: model.PartTypeUnspecified})

	r.NoError(err)
	r.Len(result, 3)
	r.Equal(uuids[0], result[0].UUID)
	r.Equal(uuids[1], result[1].UUID)
	r.Equal(uuids[2], result[2].UUID)
}
