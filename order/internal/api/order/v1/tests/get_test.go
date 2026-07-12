package tests

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestGetOrder_Success() {
	orderID := uuid.New()
	partID1 := uuid.New()
	partID2 := uuid.New()

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

	s.orderService.EXPECT().Get(s.ctx, orderID).Return(order, nil)

	params := orderv1.GetOrderParams{OrderUUID: orderID}
	resp, err := s.api.GetOrder(s.ctx, params)
	s.NoError(err)

	orderDto, ok := resp.(*orderv1.OrderDto)
	s.True(ok)
	s.Equal(orderID, orderDto.OrderUUID)
	s.Equal(partID1, orderDto.HullUUID)
	s.Equal(partID2, orderDto.EngineUUID)
	s.Equal(order.TotalPrice(), orderDto.TotalPrice)
	s.Equal(orderv1.OrderStatusPENDINGPAYMENT, orderDto.Status)
	s.True(orderDto.ShieldUUID.Null)
	s.True(orderDto.WeaponUUID.Null)
	s.True(orderDto.TransactionUUID.Null)
	s.True(orderDto.PaymentMethod.Null)
}

func (s *ServiceSuite) TestGetOrder_WithOptionalPartsAndPayment() {
	orderID := uuid.New()
	partID1 := uuid.New()
	partID2 := uuid.New()
	shieldID := uuid.New()
	weaponID := uuid.New()
	transactionID := uuid.New()
	paymentMethod := model.PaymentMethodCard

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
			{
				PartUUID: shieldID,
				PartType: model.PartTypeShield,
				Price:    500,
			},
			{
				PartUUID: weaponID,
				PartType: model.PartTypeWeapon,
				Price:    1500,
			},
		},
		TransactionUUID: &transactionID,
		PaymentMethod:   &paymentMethod,
		Status:          model.OrderStatusPaid,
		CreatedAt:       time.Now(),
	}

	s.orderService.EXPECT().Get(s.ctx, orderID).Return(order, nil)

	params := orderv1.GetOrderParams{OrderUUID: orderID}
	resp, err := s.api.GetOrder(s.ctx, params)
	s.NoError(err)

	orderDto, ok := resp.(*orderv1.OrderDto)
	s.True(ok)
	s.Equal(orderID, orderDto.OrderUUID)
	s.Equal(partID1, orderDto.HullUUID)
	s.Equal(partID2, orderDto.EngineUUID)
	s.Equal(order.TotalPrice(), orderDto.TotalPrice)
	s.Equal(orderv1.OrderStatusPAID, orderDto.Status)
	s.Equal(shieldID, orderDto.ShieldUUID.Value)
	s.Equal(weaponID, orderDto.WeaponUUID.Value)
	s.Equal(transactionID, orderDto.TransactionUUID.Value)
	s.Equal(orderv1.PaymentMethod(paymentMethod), orderDto.PaymentMethod.Value)
}

func (s *ServiceSuite) TestGetOrder_InvalidUUID() {
	params := orderv1.GetOrderParams{OrderUUID: uuid.Nil}
	resp, err := s.api.GetOrder(s.ctx, params)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *ServiceSuite) TestGetOrder_ServiceError() {
	orderID := uuid.New()
	getErr := gofakeit.Error()

	s.orderService.EXPECT().Get(s.ctx, orderID).Return(model.Order{}, getErr)

	params := orderv1.GetOrderParams{OrderUUID: orderID}
	resp, err := s.api.GetOrder(s.ctx, params)
	s.Error(err)
	s.ErrorIs(err, getErr)
	s.Nil(resp)
}
