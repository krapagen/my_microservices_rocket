package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func (s *ServiceSuite) TestCancelOrder_Success() {
	orderID := uuid.New()

	s.orderService.EXPECT().Cancel(s.ctx, orderID).Return(nil)

	params := orderv1.CancelOrderParams{OrderUUID: orderID}
	resp, err := s.api.CancelOrder(s.ctx, params)
	s.NoError(err)

	_, ok := resp.(*orderv1.CancelOrderResponse)
	s.True(ok)
}

func (s *ServiceSuite) TestCancelOrder_InvalidUUID() {
	params := orderv1.CancelOrderParams{OrderUUID: uuid.Nil}
	resp, err := s.api.CancelOrder(s.ctx, params)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *ServiceSuite) TestCancelOrder_ServiceError() {
	orderID := uuid.New()
	cancelErr := gofakeit.Error()

	s.orderService.EXPECT().Cancel(s.ctx, orderID).Return(cancelErr)

	params := orderv1.CancelOrderParams{OrderUUID: orderID}
	resp, err := s.api.CancelOrder(s.ctx, params)
	s.Error(err)
	s.ErrorIs(err, cancelErr)
	s.Nil(resp)
}
