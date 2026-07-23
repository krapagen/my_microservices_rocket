package test

import (
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestReserveParts_Success() {
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	req := &inventoryv1.ReservePartsRequest{Uuids: []string{uuid1, uuid2}}

	s.partService.EXPECT().Reserve(s.ctx, []uuid.UUID{uuid.MustParse(uuid1), uuid.MustParse(uuid2)}).Return(nil)

	resp, err := s.api.ReserveParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestReserveParts_EmptyList() {
	req := &inventoryv1.ReservePartsRequest{Uuids: []string{}}

	s.partService.EXPECT().Reserve(s.ctx, []uuid.UUID{}).Return(nil)

	resp, err := s.api.ReserveParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestReserveParts_InvalidUUID() {
	req := &inventoryv1.ReservePartsRequest{Uuids: []string{"not-a-uuid"}}

	resp, err := s.api.ReserveParts(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
	s.partService.AssertNotCalled(s.T(), "Reserve")
}

func (s *APISuite) TestReserveParts_ServiceError() {
	uuid1 := uuid.New().String()
	req := &inventoryv1.ReservePartsRequest{Uuids: []string{uuid1}}

	s.partService.EXPECT().Reserve(s.ctx, []uuid.UUID{uuid.MustParse(uuid1)}).Return(errs.ErrOutOfStock)

	resp, err := s.api.ReserveParts(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOutOfStock)
	s.Nil(resp)
}
