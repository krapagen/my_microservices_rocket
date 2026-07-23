package test

import (
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestReleaseParts_Success() {
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	req := &inventoryv1.ReleasePartsRequest{Uuids: []string{uuid1, uuid2}}

	s.partService.EXPECT().Release(s.ctx, []uuid.UUID{uuid.MustParse(uuid1), uuid.MustParse(uuid2)}).Return(nil)

	resp, err := s.api.ReleaseParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestReleaseParts_EmptyList() {
	req := &inventoryv1.ReleasePartsRequest{Uuids: []string{}}

	s.partService.EXPECT().Release(s.ctx, []uuid.UUID{}).Return(nil)

	resp, err := s.api.ReleaseParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestReleaseParts_InvalidUUID() {
	req := &inventoryv1.ReleasePartsRequest{Uuids: []string{"not-a-uuid"}}

	resp, err := s.api.ReleaseParts(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
	s.partService.AssertNotCalled(s.T(), "Release")
}

func (s *APISuite) TestReleaseParts_ServiceError() {
	uuid1 := uuid.New().String()
	req := &inventoryv1.ReleasePartsRequest{Uuids: []string{uuid1}}

	s.partService.EXPECT().Release(s.ctx, []uuid.UUID{uuid.MustParse(uuid1)}).Return(errs.ErrNothingToRelease)

	resp, err := s.api.ReleaseParts(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrNothingToRelease)
	s.Nil(resp)
}
