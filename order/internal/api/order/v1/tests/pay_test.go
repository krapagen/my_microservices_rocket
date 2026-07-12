package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	orderID := uuid.New()
	transactionID := uuid.New()

	s.orderService.EXPECT().Pay(s.ctx, orderID, model.PaymentMethod(orderv1.PaymentMethodCARD)).Return(transactionID, nil)

	req := &orderv1.PayOrderRequest{PaymentMethod: orderv1.PaymentMethodCARD}
	params := orderv1.PayOrderParams{OrderUUID: orderID}

	resp, err := s.api.PayOrder(s.ctx, req, params)
	s.NoError(err)

	payResp, ok := resp.(*orderv1.PayOrderResponse)
	s.True(ok)
	s.Equal(transactionID, payResp.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrder_InvalidUUID() {
	req := &orderv1.PayOrderRequest{PaymentMethod: orderv1.PaymentMethodCARD}
	params := orderv1.PayOrderParams{OrderUUID: uuid.Nil}

	resp, err := s.api.PayOrder(s.ctx, req, params)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *ServiceSuite) TestPayOrder_ServiceError() {
	orderID := uuid.New()
	payErr := gofakeit.Error()

	s.orderService.EXPECT().Pay(s.ctx, orderID, model.PaymentMethod(orderv1.PaymentMethodCARD)).Return(uuid.Nil, payErr)

	req := &orderv1.PayOrderRequest{PaymentMethod: orderv1.PaymentMethodCARD}
	params := orderv1.PayOrderParams{OrderUUID: orderID}

	resp, err := s.api.PayOrder(s.ctx, req, params)
	s.Error(err)
	s.ErrorIs(err, payErr)
	s.Nil(resp)
}
