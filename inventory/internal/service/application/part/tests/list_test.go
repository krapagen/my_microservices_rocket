package tests

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *ServiceSuite) TestList_Success() {
	var (
		filter = input.PartFilter{
			PartType: model.PartTypeHull,
		}
		expected = []model.Part{
			model.RestorePart(
				uuid.New(),
				gofakeit.ProductName(),
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeHull,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
		}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Equal(model.PartTypeHull, result[0].PartType())
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
			model.RestorePart(
				uuid3,
				gofakeit.ProductName(),
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeEngine,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
			model.RestorePart(
				uuid1,
				gofakeit.ProductName(),
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeHull,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
			model.RestorePart(
				uuid2,
				gofakeit.ProductName(),
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeShield,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
		}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Len(result, 3)

	s.Equal(filter.UUIDs[0], result[0].UUID()) // uuid3 (первый в запросе)
	s.Equal(filter.UUIDs[1], result[1].UUID()) // uuid1 (второй в запросе)
	s.Equal(filter.UUIDs[2], result[2].UUID()) // uuid2 (третий в запросе)
	// assert.Equal(s.T(), expected, result)
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
			model.RestorePart(
				uuid1,
				"C",
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeEngine,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
			model.RestorePart(
				uuid2,
				"B",
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeEngine,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
			model.RestorePart(
				uuid3,
				"A",
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeEngine,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
		}
	)

	s.partRepository.EXPECT().List(s.ctx, filter).Return(expected, nil)

	result, err := s.service.List(s.ctx, filter)

	s.NoError(err)
	s.Len(result, 3)

	s.Equal(uuid1, result[0].UUID()) // uuid1 (первый в запросе)
	s.Equal(uuid2, result[1].UUID()) // uuid2 (второй в запросе)
	s.Equal(uuid3, result[2].UUID()) // uuid3 (третий в запросе)
	// assert.Equal(s.T(), expected, result)
}

func (s *ServiceSuite) TestList_EmptyFilter() {
	var (
		filter   = input.PartFilter{}
		expected = []model.Part{
			model.RestorePart(
				uuid.New(),
				"Engine Part", // Начинается с E
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeEngine,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
			model.RestorePart(
				uuid.New(),
				"Hull Part", // Начинается с H (после E)
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeHull,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
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
			model.RestorePart(
				validUUID,
				gofakeit.ProductName(),
				gofakeit.LoremIpsumSentence(5),
				model.PartTypeShield,
				int64(gofakeit.Price(100, 100000)),
				int(gofakeit.Number(1, 100)),
				0,
				model.PartProperties{},
				time.Now().UTC(),
			),
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
