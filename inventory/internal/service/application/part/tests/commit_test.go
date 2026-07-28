package tests

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *ServiceSuite) TestCommit_Success() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 3, model.PartProperties{}, time.Now().UTC())
	filter := input.CommitFilter{UUIDs: []uuid.UUID{part.UUID()}}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	s.partRepository.EXPECT().ListForUpdate(s.ctx, input.PartFilter{UUIDs: filter.UUIDs}).Return([]model.Part{part}, nil)
	s.partRepository.EXPECT().Commit(s.ctx, filter).Return(nil)

	err := s.service.Commit(s.ctx, filter)
	s.NoError(err)
}

func (s *ServiceSuite) TestCommit_PartNotFound() {
	filter := input.CommitFilter{UUIDs: []uuid.UUID{uuid.New()}}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	s.partRepository.EXPECT().ListForUpdate(s.ctx, input.PartFilter{UUIDs: filter.UUIDs}).Return(nil, errs.ErrPartNotFound)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestCommit_OutOfStock() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 0, 1, model.PartProperties{}, time.Now().UTC())
	filter := input.CommitFilter{UUIDs: []uuid.UUID{part.UUID()}}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	s.partRepository.EXPECT().ListForUpdate(s.ctx, input.PartFilter{UUIDs: filter.UUIDs}).Return([]model.Part{part}, nil)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, errs.ErrCommitParts)
	s.partRepository.AssertNotCalled(s.T(), "Commit")
}

func (s *ServiceSuite) TestCommit_NotReserved() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 0, model.PartProperties{}, time.Now().UTC())
	filter := input.CommitFilter{UUIDs: []uuid.UUID{part.UUID()}}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	s.partRepository.EXPECT().ListForUpdate(s.ctx, input.PartFilter{UUIDs: filter.UUIDs}).Return([]model.Part{part}, nil)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, errs.ErrCommitParts)
	s.partRepository.AssertNotCalled(s.T(), "Commit")
}

func (s *ServiceSuite) TestCommit_RepositoryCommitError() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 2, model.PartProperties{}, time.Now().UTC())
	filter := input.CommitFilter{UUIDs: []uuid.UUID{part.UUID()}}
	repoErr := gofakeit.Error()

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	s.partRepository.EXPECT().ListForUpdate(s.ctx, input.PartFilter{UUIDs: filter.UUIDs}).Return([]model.Part{part}, nil)
	s.partRepository.EXPECT().Commit(s.ctx, filter).Return(repoErr)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCommit_RepositoryCommitPartsError() {
	part := model.RestorePart(uuid.New(), "Hull", "", model.PartTypeHull, 1000, 10, 2, model.PartProperties{}, time.Now().UTC())
	filter := input.CommitFilter{UUIDs: []uuid.UUID{part.UUID()}}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
	s.partRepository.EXPECT().ListForUpdate(s.ctx, input.PartFilter{UUIDs: filter.UUIDs}).Return([]model.Part{part}, nil)
	s.partRepository.EXPECT().Commit(s.ctx, filter).Return(errs.ErrCommitParts)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, errs.ErrCommitParts)
}
