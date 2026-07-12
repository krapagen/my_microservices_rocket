package test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/payment/internal/errors"
	"github.com/krapagen/my_microservices_rocket/payment/internal/model"
	"github.com/krapagen/my_microservices_rocket/payment/internal/service/input"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrder_Success() {
	orderID := uuid.New()
	transactionID := uuid.New()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	expectedInput := input.PayOrderInput{
		OrderUUID:     orderID,
		PaymentMethod: model.PaymentMethodCard,
	}

	s.paymentService.EXPECT().Pay(s.ctx, expectedInput).Return(transactionID, nil)

	resp, err := s.api.PayOrder(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(transactionID.String(), resp.TransactionUuid)
}

func (s *APISuite) TestPayOrder_InvalidOrderUUID() {
	req := &paymentv1.PayOrderRequest{
		OrderUuid:     "not-a-uuid",
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	resp, err := s.api.PayOrder(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidOrderUUID)
	s.Nil(resp)
}

func (s *APISuite) TestPayOrder_InvalidPaymentMethod() {
	orderID := uuid.New()
	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	}

	resp, err := s.api.PayOrder(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidPaymentMethod)
	s.Nil(resp)
}

func (s *APISuite) TestPayOrder_AllPaymentMethods() {
	orderID := uuid.New()

	cases := []struct {
		protoMethod paymentv1.PaymentMethod
		modelMethod model.PaymentMethod
	}{
		{paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, model.PaymentMethodSBP},
		{paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, model.PaymentMethodCreditCard},
		{paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, model.PaymentMethodInvestorMoney},
	}

	for _, tc := range cases {
		s.Run(string(tc.modelMethod), func() {
			req := &paymentv1.PayOrderRequest{
				OrderUuid:     orderID.String(),
				PaymentMethod: tc.protoMethod,
			}
			expectedInput := input.PayOrderInput{
				OrderUUID:     orderID,
				PaymentMethod: tc.modelMethod,
			}
			transactionID := uuid.New()

			s.paymentService.EXPECT().Pay(s.ctx, expectedInput).Return(transactionID, nil)

			resp, err := s.api.PayOrder(s.ctx, req)
			s.NoError(err)
			s.NotNil(resp)
			s.Equal(transactionID.String(), resp.TransactionUuid)
		})
	}
}

func (s *APISuite) TestPayOrder_ServiceError() {
	orderID := uuid.New()
	payErr := gofakeit.Error()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	expectedInput := input.PayOrderInput{
		OrderUUID:     orderID,
		PaymentMethod: model.PaymentMethodCard,
	}

	s.paymentService.EXPECT().Pay(s.ctx, expectedInput).Return(uuid.Nil, payErr)

	resp, err := s.api.PayOrder(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, payErr)
	s.Nil(resp)
}
