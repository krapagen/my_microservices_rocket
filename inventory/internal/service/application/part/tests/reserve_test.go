package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func (s *ServiceSuite) TestReserve_Success() {
	hull := newFakeHull(100)
	engine := newFakeEngine(model.EngineClassC, 30)
	uuids := []uuid.UUID{hull.UUID(), engine.UUID()}

	s.partRepository.EXPECT().ReserveAtomic(s.ctx, uuids).Return(nil)

	err := s.service.Reserve(s.ctx, uuids)
	s.NoError(err)
}

func (s *ServiceSuite) TestReserve_EmptyUUIDs() {
	err := s.service.Reserve(s.ctx, []uuid.UUID{})
	s.NoError(err)
	s.txManager.AssertNotCalled(s.T(), "Do")
	s.partRepository.AssertNotCalled(s.T(), "ReserveAtomic")
}

func (s *ServiceSuite) TestReserve_PartNotFound() {
	uuids := []uuid.UUID{uuid.New()}
	s.partRepository.EXPECT().ReserveAtomic(s.ctx, uuids).Return(errs.ErrPartNotFound)

	err := s.service.Reserve(s.ctx, uuids)
	s.Error(err)
	s.ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestReserve_OutOfStock() {
	uuids := []uuid.UUID{uuid.New()}
	s.partRepository.EXPECT().ReserveAtomic(s.ctx, uuids).Return(errs.ErrOutOfStock)

	err := s.service.Reserve(s.ctx, uuids)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOutOfStock)
}

func (s *ServiceSuite) TestReserve_UpdateError() {
	uuids := []uuid.UUID{uuid.New()}
	updateErr := gofakeit.Error()

	s.partRepository.EXPECT().ReserveAtomic(s.ctx, uuids).Return(updateErr)

	err := s.service.Reserve(s.ctx, uuids)
	s.ErrorIs(err, updateErr)
}
