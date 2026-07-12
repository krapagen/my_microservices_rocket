package test

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPart_Success() {
	var (
		partUUID = uuid.New()
		part     = model.Part{
			UUID:          partUUID,
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.LoremIpsumSentence(5),
			Price:         int64(gofakeit.Price(100, 100000)),
			PartType:      model.PartTypeEngine,
			StockQuantity: int64(gofakeit.Number(1, 100)),
			CreatedAt:     time.Now().UTC(),
		}
		req = &inventoryv1.GetPartRequest{Uuid: partUUID.String()}
	)

	s.partService.EXPECT().Get(s.ctx, partUUID).Return(part, nil)

	resp, err := s.api.GetPart(s.ctx, req)

	s.NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.GetPart())
	s.Equal(partUUID.String(), resp.GetPart().GetUuid())
	s.Equal(part.Name, resp.GetPart().GetName())
	s.Equal(part.Description, resp.GetPart().GetDescription())
	s.Equal(part.Price, resp.GetPart().GetPrice())
	s.Equal(inventoryv1.PartType_PART_TYPE_ENGINE, resp.GetPart().GetPartType())
	s.Equal(part.StockQuantity, resp.GetPart().GetStockQuantity())
	s.Equal(part.CreatedAt, resp.GetPart().GetCreatedAt().AsTime())
}

func (s *APISuite) TestGetPart_InvalidUUID() {
	req := &inventoryv1.GetPartRequest{Uuid: "not-a-uuid"}

	resp, err := s.api.GetPart(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
	s.partService.AssertNotCalled(s.T(), "Get")
}

func (s *APISuite) TestGetPart_NotFound() {
	var (
		partUUID = uuid.New()
		req      = &inventoryv1.GetPartRequest{Uuid: partUUID.String()}
	)

	s.partService.EXPECT().Get(s.ctx, partUUID).Return(model.Part{}, errs.ErrPartNotFound)

	resp, err := s.api.GetPart(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrPartNotFound)
	s.Nil(resp)
}

func (s *APISuite) TestGetPart_ServiceError() {
	var (
		partUUID = uuid.New()
		svcErr   = gofakeit.Error()
		req      = &inventoryv1.GetPartRequest{Uuid: partUUID.String()}
	)

	s.partService.EXPECT().Get(s.ctx, partUUID).Return(model.Part{}, svcErr)

	resp, err := s.api.GetPart(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}
