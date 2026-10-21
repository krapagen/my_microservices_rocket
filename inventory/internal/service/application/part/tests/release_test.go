package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

func (s *ServiceSuite) TestRelease_Success() {
	uuids := []uuid.UUID{uuid.New()}

	s.partRepository.EXPECT().ReleaseAtomic(s.ctx, uuids).Return(nil)

	err := s.service.Release(s.ctx, uuids)
	s.NoError(err)
}

func (s *ServiceSuite) TestRelease_EmptyUUIDs() {
	err := s.service.Release(s.ctx, []uuid.UUID{})
	s.NoError(err)
	s.txManager.AssertNotCalled(s.T(), "Do")
	s.partRepository.AssertNotCalled(s.T(), "ReleaseAtomic")
}

func (s *ServiceSuite) TestRelease_PartNotFound() {
	uuids := []uuid.UUID{uuid.New()}
	s.partRepository.EXPECT().ReleaseAtomic(s.ctx, uuids).Return(errs.ErrPartNotFound)

	err := s.service.Release(s.ctx, uuids)
	s.ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestRelease_NothingToRelease() {
	uuids := []uuid.UUID{uuid.New()}
	s.partRepository.EXPECT().ReleaseAtomic(s.ctx, uuids).Return(errs.ErrNothingToRelease)

	err := s.service.Release(s.ctx, uuids)
	s.Error(err)
	s.ErrorIs(err, errs.ErrNothingToRelease)
}

func (s *ServiceSuite) TestRelease_UpdateError() {
	uuids := []uuid.UUID{uuid.New()}
	updateErr := gofakeit.Error()

	s.partRepository.EXPECT().ReleaseAtomic(s.ctx, uuids).Return(updateErr)

	err := s.service.Release(s.ctx, uuids)
	s.ErrorIs(err, updateErr)
}
