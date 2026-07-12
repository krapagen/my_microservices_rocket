package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestPay_Success() {
	orderID := uuid.New()
	partID := uuid.New()
	transactionID := uuid.New()
	paymentMethod := model.PaymentMethodCard

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
	s.orderPaymentClient.EXPECT().PayOrder(s.ctx, orderID, paymentMethod).Return(transactionID, nil)
	s.orderRepository.EXPECT().Update(s.ctx, mock.MatchedBy(func(order model.Order) bool {
		return order.UUID == orderID &&
			order.Status == model.OrderStatusPaid &&
			order.TransactionUUID != nil &&
			*order.TransactionUUID == transactionID &&
			order.PaymentMethod != nil &&
			*order.PaymentMethod == paymentMethod
	})).Return(nil)

	result, err := s.service.Pay(s.ctx, orderID, paymentMethod)
	s.NoError(err)
	s.Equal(transactionID, result)
}

func (s *ServiceSuite) TestPay_NotFound() {
	orderID := uuid.New()
	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(model.Order{}, errs.ErrOrderNotFound)

	result, err := s.service.Pay(s.ctx, orderID, model.PaymentMethodCard)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderNotFound)
	s.Equal(uuid.Nil, result)
}

func (s *ServiceSuite) TestPay_GetRepoError() {
	orderID := uuid.New()
	repoErr := gofakeit.Error()
	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(model.Order{}, repoErr)

	result, err := s.service.Pay(s.ctx, orderID, model.PaymentMethodCard)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Equal(uuid.Nil, result)
}

func (s *ServiceSuite) TestPay_AlreadyPaid() {
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

	result, err := s.service.Pay(s.ctx, orderID, model.PaymentMethodCard)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderAlreadyPaid)
	s.Equal(uuid.Nil, result)
}

func (s *ServiceSuite) TestPay_Cancelled() {
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

	result, err := s.service.Pay(s.ctx, orderID, model.PaymentMethodCard)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderCancelled)
	s.Equal(uuid.Nil, result)
}

func (s *ServiceSuite) TestPay_UnknownStatus() {
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

	result, err := s.service.Pay(s.ctx, orderID, model.PaymentMethodCard)
	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderStatusIncorrect)
	s.Equal(uuid.Nil, result)
}

func (s *ServiceSuite) TestPay_PaymentError() {
	orderID := uuid.New()
	partID := uuid.New()
	paymentMethod := model.PaymentMethodCard
	paymentErr := gofakeit.Error()

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
	s.orderPaymentClient.EXPECT().PayOrder(s.ctx, orderID, paymentMethod).Return(uuid.Nil, paymentErr)

	result, err := s.service.Pay(s.ctx, orderID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, paymentErr)
	s.Equal(uuid.Nil, result)
}

func (s *ServiceSuite) TestPay_UpdateError() {
	orderID := uuid.New()
	partID := uuid.New()
	transactionID := uuid.New()
	paymentMethod := model.PaymentMethodCard
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
	s.orderPaymentClient.EXPECT().PayOrder(s.ctx, orderID, paymentMethod).Return(transactionID, nil)
	s.orderRepository.EXPECT().Update(s.ctx, mock.Anything).Return(updateErr)

	result, err := s.service.Pay(s.ctx, orderID, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, updateErr)
	s.Equal(uuid.Nil, result)
}
