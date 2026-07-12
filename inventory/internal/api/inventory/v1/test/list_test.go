package test

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func newFakePart(partType model.PartType) model.Part {
	return model.Part{
		UUID:          uuid.New(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.LoremIpsumSentence(5),
		Price:         int64(gofakeit.Price(100, 100000)),
		PartType:      partType,
		StockQuantity: int64(gofakeit.Number(1, 100)),
		CreatedAt:     time.Now().UTC(),
	}
}

func (s *APISuite) TestListParts_SuccessByType() {
	var (
		parts = []model.Part{
			newFakePart(model.PartTypeHull),
			newFakePart(model.PartTypeHull),
		}
		filter = input.PartFilter{
			UUIDs:    []uuid.UUID{},
			PartType: model.PartTypeHull,
		}
		req = &inventoryv1.ListPartsRequest{
			PartType: inventoryv1.PartType_PART_TYPE_HULL,
		}
	)

	s.partService.EXPECT().List(s.ctx, filter).Return(parts, nil)

	resp, err := s.api.ListParts(s.ctx, req)

	s.NoError(err)
	s.Require().NotNil(resp)
	s.Require().Len(resp.GetParts(), 2)
	for i, p := range resp.GetParts() {
		s.Equal(parts[i].UUID.String(), p.GetUuid())
		s.Equal(parts[i].Name, p.GetName())
		s.Equal(parts[i].Description, p.GetDescription())
		s.Equal(parts[i].Price, p.GetPrice())
		s.Equal(inventoryv1.PartType_PART_TYPE_HULL, p.GetPartType())
		s.Equal(parts[i].StockQuantity, p.GetStockQuantity())
		s.Equal(parts[i].CreatedAt, p.GetCreatedAt().AsTime())
	}
}

func (s *APISuite) TestListParts_SuccessByUUIDs() {
	var (
		uuid1  = uuid.New()
		uuid2  = uuid.New()
		parts  = []model.Part{newFakePart(model.PartTypeEngine), newFakePart(model.PartTypeShield)}
		filter = input.PartFilter{
			UUIDs:    []uuid.UUID{uuid1, uuid2},
			PartType: model.PartTypeUnspecified,
		}
		req = &inventoryv1.ListPartsRequest{
			Uuids: []string{uuid1.String(), uuid2.String()},
		}
	)

	s.partService.EXPECT().List(s.ctx, filter).Return(parts, nil)

	resp, err := s.api.ListParts(s.ctx, req)

	s.NoError(err)
	s.Require().NotNil(resp)
	s.Len(resp.GetParts(), 2)
}

func (s *APISuite) TestListParts_EmptyResult() {
	var (
		filter = input.PartFilter{
			UUIDs:    []uuid.UUID{},
			PartType: model.PartTypeWeapon,
		}
		req = &inventoryv1.ListPartsRequest{
			PartType: inventoryv1.PartType_PART_TYPE_WEAPON,
		}
	)

	s.partService.EXPECT().List(s.ctx, filter).Return([]model.Part{}, nil)

	resp, err := s.api.ListParts(s.ctx, req)

	s.NoError(err)
	s.Require().NotNil(resp)
	s.Empty(resp.GetParts())
}

func (s *APISuite) TestListParts_InvalidUUID() {
	req := &inventoryv1.ListPartsRequest{
		Uuids: []string{uuid.New().String(), "not-a-uuid"},
	}

	resp, err := s.api.ListParts(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
	s.partService.AssertNotCalled(s.T(), "List")
}

func (s *APISuite) TestListParts_ServiceError() {
	var (
		svcErr = gofakeit.Error()
		filter = input.PartFilter{
			UUIDs:    []uuid.UUID{},
			PartType: model.PartTypeUnspecified,
		}
		req = &inventoryv1.ListPartsRequest{}
	)

	s.partService.EXPECT().List(s.ctx, filter).Return(nil, svcErr)

	resp, err := s.api.ListParts(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}
