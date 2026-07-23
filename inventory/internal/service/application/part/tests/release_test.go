package tests

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *ServiceSuite) TestRelease_Success() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 3, model.PartProperties{}, time.Now().UTC())
	uuids := []uuid.UUID{part.UUID()}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: uuids}).Return([]model.Part{part}, nil)
	s.partRepository.EXPECT().UpdateReservedBatch(s.ctx, mock.Anything).Return(nil)

	err := s.service.Release(s.ctx, uuids)
	s.NoError(err)
}

func (s *ServiceSuite) TestRelease_EmptyUUIDs() {
	err := s.service.Release(s.ctx, []uuid.UUID{})
	s.NoError(err)
	s.partRepository.AssertNotCalled(s.T(), "List")
	s.partRepository.AssertNotCalled(s.T(), "UpdateReservedBatch")
}

func (s *ServiceSuite) TestRelease_PartNotFound() {
	uuids := []uuid.UUID{uuid.New()}
	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: uuids}).Return(nil, errs.ErrPartNotFound)

	err := s.service.Release(s.ctx, uuids)
	s.ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestRelease_NothingToRelease() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 0, model.PartProperties{}, time.Now().UTC())
	uuids := []uuid.UUID{part.UUID()}
	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: uuids}).Return([]model.Part{part}, nil)

	err := s.service.Release(s.ctx, uuids)
	s.Error(err)
	s.ErrorIs(err, errs.ErrNothingToRelease)
}

func (s *ServiceSuite) TestRelease_UpdateError() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 3, model.PartProperties{}, time.Now().UTC())
	uuids := []uuid.UUID{part.UUID()}
	updateErr := gofakeit.Error()

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: uuids}).Return([]model.Part{part}, nil)
	s.partRepository.EXPECT().UpdateReservedBatch(s.ctx, mock.Anything).Return(updateErr)

	err := s.service.Release(s.ctx, uuids)
	s.ErrorIs(err, updateErr)
}
