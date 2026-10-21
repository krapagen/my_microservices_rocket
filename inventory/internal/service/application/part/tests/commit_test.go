package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *ServiceSuite) TestCommit_Success() {
	filter := input.CommitFilter{UUIDs: []uuid.UUID{uuid.New()}}

	s.partRepository.EXPECT().Commit(s.ctx, filter).Return(nil)

	err := s.service.Commit(s.ctx, filter)
	s.NoError(err)
}

func (s *ServiceSuite) TestCommit_RepositoryCommitError() {
	filter := input.CommitFilter{UUIDs: []uuid.UUID{uuid.New()}}
	repoErr := gofakeit.Error()

	s.partRepository.EXPECT().Commit(s.ctx, filter).Return(repoErr)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCommit_RepositoryCommitPartsError() {
	filter := input.CommitFilter{UUIDs: []uuid.UUID{uuid.New()}}

	s.partRepository.EXPECT().Commit(s.ctx, filter).Return(errs.ErrCommitParts)

	err := s.service.Commit(s.ctx, filter)
	s.ErrorIs(err, errs.ErrCommitParts)
}
