package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/input"
)

func (s *ServiceSuite) TestCreate_Success() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	partUUIDs := []uuid.UUID{partID1, partID2}
	modelParts := []model.Part{
		{
			UUID:          partID1,
			Name:          "Hull Part",
			PartType:      model.PartTypeHull,
			Price:         1000,
			StockQuantity: 10,
		},
		{
			UUID:          partID2,
			Name:          "Engine Part",
			PartType:      model.PartTypeEngine,
			Price:         2000,
			StockQuantity: 5,
		},
	}
	expectedOrder := model.Order{
		Items: []model.OrderItem{
			{
				PartUUID: partID1,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
			{
				PartUUID: partID2,
				PartType: model.PartTypeEngine,
				Price:    2000,
			},
		},
		Status: model.OrderStatusPendingPayment,
	}

	s.orderInventoryClient.EXPECT().ListParts(s.ctx, partUUIDs).Return(modelParts, nil)
	s.orderRepository.EXPECT().Create(s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID != uuid.Nil &&
			!order.CreatedAt.IsZero() &&
			order.Status == expectedOrder.Status &&
			len(order.Items) == len(expectedOrder.Items) &&
			order.Items[0] == expectedOrder.Items[0] &&
			order.Items[1] == expectedOrder.Items[1]
	})).Return(nil)

	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	})
	s.NoError(err)
	s.Equal(expectedOrder.Status, result.Status)
	s.Equal(expectedOrder.Items, result.Items)
	s.NotEqual(uuid.Nil, result.UUID)
	s.False(result.CreatedAt.IsZero())
}

func (s *ServiceSuite) TestCreate_MissingRequiredParts() {
	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   uuid.Nil,
		EngineUUID: uuid.New(),
	})
	s.Error(err)
	s.ErrorIs(err, errs.ErrMissingRequiredParts)
	s.Equal(model.Order{}, result)
}

func (s *ServiceSuite) TestCreate_InventoryError() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	partUUIDs := []uuid.UUID{partID1, partID2}
	invErr := gofakeit.Error()

	s.orderInventoryClient.EXPECT().ListParts(s.ctx, partUUIDs).Return(nil, invErr)

	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	})
	s.Error(err)
	s.ErrorIs(err, invErr)
	s.Equal(model.Order{}, result)
}

func (s *ServiceSuite) TestCreate_OutOfStock() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	partUUIDs := []uuid.UUID{partID1, partID2}
	modelParts := []model.Part{
		{
			UUID:          partID1,
			Name:          "Hull Part",
			PartType:      model.PartTypeHull,
			Price:         1000,
			StockQuantity: 0,
		},
		{
			UUID:          partID2,
			Name:          "Engine Part",
			PartType:      model.PartTypeEngine,
			Price:         2000,
			StockQuantity: 5,
		},
	}

	s.orderInventoryClient.EXPECT().ListParts(s.ctx, partUUIDs).Return(modelParts, nil)

	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	})
	s.Error(err)
	s.ErrorIs(err, errs.ErrOutOfStock)
	s.Equal(model.Order{}, result)
}

func (s *ServiceSuite) TestCreate_EmptyParts() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	partUUIDs := []uuid.UUID{partID1, partID2}

	s.orderInventoryClient.EXPECT().ListParts(s.ctx, partUUIDs).Return([]model.Part{}, nil)

	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	})
	s.Error(err)
	s.ErrorIs(err, errs.ErrMissingRequiredParts)
	s.Equal(model.Order{}, result)
}

func (s *ServiceSuite) TestCreate_WithOptionalParts() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	shieldID := uuid.New()
	weaponID := uuid.New()
	partUUIDs := []uuid.UUID{partID1, partID2, shieldID, weaponID}
	modelParts := []model.Part{
		{
			UUID:          partID1,
			Name:          "Hull Part",
			PartType:      model.PartTypeHull,
			Price:         1000,
			StockQuantity: 10,
		},
		{
			UUID:          partID2,
			Name:          "Engine Part",
			PartType:      model.PartTypeEngine,
			Price:         2000,
			StockQuantity: 5,
		},
		{
			UUID:          shieldID,
			Name:          "Shield Part",
			PartType:      model.PartTypeShield,
			Price:         500,
			StockQuantity: 3,
		},
		{
			UUID:          weaponID,
			Name:          "Weapon Part",
			PartType:      model.PartTypeWeapon,
			Price:         1500,
			StockQuantity: 2,
		},
	}

	s.orderInventoryClient.EXPECT().ListParts(s.ctx, partUUIDs).Return(modelParts, nil)
	s.orderRepository.EXPECT().Create(s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID != uuid.Nil &&
			!order.CreatedAt.IsZero() &&
			order.Status == model.OrderStatusPendingPayment &&
			len(order.Items) == 4
	})).Return(nil)

	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
		ShieldUUID: &shieldID,
		WeaponUUID: &weaponID,
	})
	s.NoError(err)
	s.Equal(model.OrderStatusPendingPayment, result.Status)
	s.Len(result.Items, 4)
	s.NotEqual(uuid.Nil, result.UUID)
	s.False(result.CreatedAt.IsZero())
}

func (s *ServiceSuite) TestCreate_RepositoryError() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	partUUIDs := []uuid.UUID{partID1, partID2}
	modelParts := []model.Part{
		{
			UUID:          partID1,
			Name:          "Hull Part",
			PartType:      model.PartTypeHull,
			Price:         1000,
			StockQuantity: 10,
		},
		{
			UUID:          partID2,
			Name:          "Engine Part",
			PartType:      model.PartTypeEngine,
			Price:         2000,
			StockQuantity: 5,
		},
	}
	repoErr := gofakeit.Error()

	s.orderInventoryClient.EXPECT().ListParts(s.ctx, partUUIDs).Return(modelParts, nil)
	s.orderRepository.EXPECT().Create(s.ctx, mock.Anything).Return(repoErr)

	result, err := s.service.Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	})
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Equal(model.Order{}, result)
}
