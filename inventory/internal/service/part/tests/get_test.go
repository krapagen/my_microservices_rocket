package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func (s *ServiceSuite) TestGet_Success() {
	var (
		partUUID = uuid.New()
		expected = model.Part{
			UUID:          partUUID,
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.LoremIpsumSentence(5),
			Price:         int64(gofakeit.Price(100, 100000)),
			PartType:      model.PartTypeHull,
			StockQuantity: int64(gofakeit.Number(1, 100)),
		}
	)

	s.partRepository.EXPECT().Get(s.ctx, partUUID).Return(expected, nil)

	result, err := s.service.Get(s.ctx, partUUID)

	s.NoError(err)
	s.Equal(expected, result)
	assert.Equal(s.T(), expected, result)
}

func (s *ServiceSuite) TestGet_NotFound() {
	partUUID := uuid.New()

	s.partRepository.EXPECT().Get(s.ctx, partUUID).Return(model.Part{}, errs.ErrPartNotFound)

	result, err := s.service.Get(s.ctx, partUUID)

	s.Error(err)
	s.ErrorIs(err, errs.ErrPartNotFound)
	s.Equal(model.Part{}, result)
}

func (s *ServiceSuite) TestGet_RepoError() {
	var (
		partUUID = uuid.New()
		repoErr  = gofakeit.Error()
	)

	s.partRepository.EXPECT().Get(s.ctx, partUUID).Return(model.Part{}, repoErr)

	result, err := s.service.Get(s.ctx, partUUID)

	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Equal(model.Part{}, result)
}
