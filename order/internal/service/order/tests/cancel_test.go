package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCancel_Success() {
	orderID := uuid.New()
	partID := uuid.New()

	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status: model.OrderStatusPendingPayment,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(testOrder, nil)
	s.orderRepository.EXPECT().Update(s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderID && order.Status == model.OrderStatusCancelled
	})).Return(nil)

	err := s.service.Cancel(s.ctx, orderID)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancel_NotFound() {
	orderID := uuid.New()
	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(model.Order{}, errs.ErrOrderNotFound)

	err := s.service.Cancel(s.ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancel_GetRepoError() {
	orderID := uuid.New()
	repoErr := gofakeit.Error()
	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(model.Order{}, repoErr)

	err := s.service.Cancel(s.ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCancel_AlreadyPaid() {
	orderID := uuid.New()
	partID := uuid.New()
	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status: model.OrderStatusPaid,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(testOrder, nil)

	err := s.service.Cancel(s.ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancel_AlreadyCancelled() {
	orderID := uuid.New()
	partID := uuid.New()
	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status: model.OrderStatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(testOrder, nil)

	err := s.service.Cancel(s.ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderCancelled)
}

func (s *ServiceSuite) TestCancel_UnknownStatus() {
	orderID := uuid.New()
	partID := uuid.New()
	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status: model.OrderStatus("UNKNOWN"),
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(testOrder, nil)

	err := s.service.Cancel(s.ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderStatusIncorrect)
}

func (s *ServiceSuite) TestCancel_UpdateError() {
	orderID := uuid.New()
	partID := uuid.New()
	updateErr := gofakeit.Error()

	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status: model.OrderStatusPendingPayment,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(testOrder, nil)
	s.orderRepository.EXPECT().Update(s.ctx, mock.Anything).Return(updateErr)

	err := s.service.Cancel(s.ctx, orderID)
	s.Error(err)
	s.ErrorIs(err, updateErr)
}
