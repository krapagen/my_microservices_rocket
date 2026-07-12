package tests

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/input"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	orderID := uuid.New()

	order := model.Order{
		UUID: orderID,
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
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	s.orderService.EXPECT().Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	}).Return(order, nil)

	req := &orderv1.CreateOrderRequest{
		HullUUID:   partID1,
		EngineUUID: partID2,
	}

	resp, err := s.api.CreateOrder(s.ctx, req)
	s.NoError(err)

	createResp, ok := resp.(*orderv1.CreateOrderResponse)
	s.True(ok)
	s.Equal(orderID, createResp.OrderUUID)
	s.Equal(order.TotalPrice(), createResp.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrder_MissingRequiredParts() {
	req := &orderv1.CreateOrderRequest{
		HullUUID:   uuid.Nil,
		EngineUUID: uuid.New(),
	}

	resp, err := s.api.CreateOrder(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrMissingRequiredParts)
	s.Nil(resp)
}

func (s *ServiceSuite) TestCreateOrder_WithOptionalParts() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	shieldID := uuid.New()
	weaponID := uuid.New()
	orderID := uuid.New()

	order := model.Order{
		UUID: orderID,
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
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	s.orderService.EXPECT().Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
		ShieldUUID: &shieldID,
		WeaponUUID: &weaponID,
	}).Return(order, nil)

	req := &orderv1.CreateOrderRequest{
		HullUUID:   partID1,
		EngineUUID: partID2,
		ShieldUUID: orderv1.OptNilUUID{Value: shieldID, Set: true, Null: false},
		WeaponUUID: orderv1.OptNilUUID{Value: weaponID, Set: true, Null: false},
	}

	resp, err := s.api.CreateOrder(s.ctx, req)
	s.NoError(err)

	createResp, ok := resp.(*orderv1.CreateOrderResponse)
	s.True(ok)
	s.Equal(orderID, createResp.OrderUUID)
	s.Equal(order.TotalPrice(), createResp.TotalPrice)
}

func (s *ServiceSuite) TestCreateOrder_ServiceError() {
	partID1 := uuid.New()
	partID2 := uuid.New()
	createErr := gofakeit.Error()

	s.orderService.EXPECT().Create(s.ctx, input.CreateOrderInput{
		HullUUID:   partID1,
		EngineUUID: partID2,
	}).Return(model.Order{}, createErr)

	req := &orderv1.CreateOrderRequest{
		HullUUID:   partID1,
		EngineUUID: partID2,
	}

	resp, err := s.api.CreateOrder(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, createErr)
	s.Nil(resp)
}
