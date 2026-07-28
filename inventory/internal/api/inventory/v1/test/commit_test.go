package test

import (
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestCommitParts_Success() {
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	req := &inventoryv1.CommitPartsRequest{Uuids: []string{uuid1, uuid2}}

	s.partService.EXPECT().Commit(s.ctx, input.CommitFilter{
		UUIDs: []uuid.UUID{uuid.MustParse(uuid1), uuid.MustParse(uuid2)},
	}).Return(nil)

	resp, err := s.api.CommitParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestCommitParts_EmptyList() {
	req := &inventoryv1.CommitPartsRequest{Uuids: []string{}}

	s.partService.EXPECT().Commit(s.ctx, input.CommitFilter{UUIDs: []uuid.UUID{}}).Return(nil)

	resp, err := s.api.CommitParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestCommitParts_InvalidUUID() {
	req := &inventoryv1.CommitPartsRequest{Uuids: []string{"not-a-uuid"}}

	resp, err := s.api.CommitParts(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
	s.partService.AssertNotCalled(s.T(), "Commit")
}

func (s *APISuite) TestCommitParts_ServiceError() {
	uuid1 := uuid.New().String()
	req := &inventoryv1.CommitPartsRequest{Uuids: []string{uuid1}}

	s.partService.EXPECT().Commit(s.ctx, input.CommitFilter{
		UUIDs: []uuid.UUID{uuid.MustParse(uuid1)},
	}).Return(errs.ErrCommitParts)

	resp, err := s.api.CommitParts(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrCommitParts)
	s.Nil(resp)
}
