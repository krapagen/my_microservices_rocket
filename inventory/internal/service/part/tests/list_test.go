package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/stretchr/testify/assert"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *ServiceSuite) TestList_Success() {
	var (
		filter = input.PartFilter{
			PartType: model.PartTypeHull,
		}
		expected = []model.Part{
			{
				UUID:          uuid.New(),
				Name:          gofakeit.ProductName(),
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeHull,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
		}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Equal(model.PartTypeHull, result[0].PartType)
	assert.Equal(s.T(), expected, result)
}

func (s *ServiceSuite) TestList_EmptyResult() {
	var (
		filter = input.PartFilter{
			PartType: model.PartTypeWeapon,
		}
		expected = []model.Part{}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Equal(expected, result)
	assert.Equal(s.T(), expected, result)
}

func (s *ServiceSuite) TestList_WithMultipleUUIDs() {
	var (
		uuid1  = uuid.New()
		uuid2  = uuid.New()
		uuid3  = uuid.New()
		filter = input.PartFilter{
			UUIDs: []uuid.UUID{uuid3, uuid1, uuid2}, // в определенном порядке
		}
		expected = []model.Part{
			{
				UUID:          uuid3,
				Name:          gofakeit.ProductName(),
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeEngine,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
			{
				UUID:          uuid1,
				Name:          gofakeit.ProductName(),
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeHull,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
			{
				UUID:          uuid2,
				Name:          gofakeit.ProductName(),
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeShield,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
		}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Len(result, 3)

	s.Equal(filter.UUIDs[0], result[0].UUID) // uuid3 (первый в запросе)
	s.Equal(filter.UUIDs[1], result[1].UUID) // uuid1 (второй в запросе)
	s.Equal(filter.UUIDs[2], result[2].UUID) // uuid2 (третий в запросе)
	//assert.Equal(s.T(), expected, result)
}

func (s *ServiceSuite) TestList_Sorted() {
	var (
		uuid1  = uuid.New()
		uuid2  = uuid.New()
		uuid3  = uuid.New()
		filter = input.PartFilter{
			PartType: model.PartTypeEngine,
		}
		expected = []model.Part{
			{
				UUID:          uuid3,
				Name:          "C",
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeEngine,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
			{
				UUID:          uuid1,
				Name:          "B",
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeEngine,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
			{
				UUID:          uuid2,
				Name:          "A",
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeEngine,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
		}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Len(result, 3)

	s.Equal(uuid2, result[0].UUID) // uuid3 (первый в запросе)
	s.Equal(uuid1, result[1].UUID) // uuid1 (второй в запросе)
	s.Equal(uuid3, result[2].UUID) // uuid2 (третий в запросе)
	//assert.Equal(s.T(), expected, result)
}

func (s *ServiceSuite) TestList_EmptyFilter() {
	var (
		filter   = input.PartFilter{}
		expected = []model.Part{
			{
				UUID:          uuid.New(),
				Name:          "Engine Part", // Начинается с E
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeEngine,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
			{
				UUID:          uuid.New(),
				Name:          "Hull Part", // Начинается с H (после E)
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeHull,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
		}
		dst = make([]model.Part, 2)
	)
	copy(dst, expected)
	s.partRepository.EXPECT().List(s.ctx, filter).Return(dst, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Len(result, 2)
	s.Equal(expected[0], result[0])
	s.Equal(expected[1], result[1])
}

func (s *ServiceSuite) TestList_WithEmptyUUID() {
	var (
		validUUID = uuid.New()
		filter    = input.PartFilter{
			UUIDs: []uuid.UUID{uuid.Nil, validUUID}, // один пустой UUID и один валидный
		}
		expected = []model.Part{
			{
				UUID:          validUUID,
				Name:          gofakeit.ProductName(),
				Description:   gofakeit.LoremIpsumSentence(5),
				Price:         int64(gofakeit.Price(100, 100000)),
				PartType:      model.PartTypeShield,
				StockQuantity: int64(gofakeit.Number(1, 100)),
			},
		}
		dst = make([]model.Part, 2)
	)
	copy(dst, expected)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(dst, errs.ErrPartNotFound)

	_, err := s.service.List(s.ctx, filter)

	s.ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestList_RepoError() {
	var (
		filter  = input.PartFilter{}
		repoErr = gofakeit.Error()
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return([]model.Part{}, repoErr)

	_, err := s.service.List(s.ctx, filter)

	s.Error(err)
	s.ErrorIs(err, repoErr)
}
