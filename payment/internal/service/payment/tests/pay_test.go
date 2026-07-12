package tests

import (
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/payment/internal/errors"
	"github.com/krapagen/my_microservices_rocket/payment/internal/model"
	"github.com/krapagen/my_microservices_rocket/payment/internal/service/input"
)

func (s *ServiceSuite) TestPay_Success() {
	in := input.PayOrderInput{
		OrderUUID:     uuid.New(),
		PaymentMethod: model.PaymentMethodCard,
	}

	transactionID, err := s.service.Pay(s.ctx, in)
	s.NoError(err)
	s.NotEqual(uuid.Nil, transactionID)
}

func (s *ServiceSuite) TestPay_InvalidPaymentMethod() {
	in := input.PayOrderInput{
		OrderUUID:     uuid.New(),
		PaymentMethod: model.PaymentMethodUnspecified,
	}

	transactionID, err := s.service.Pay(s.ctx, in)
	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidPaymentMethod)
	s.Equal(uuid.Nil, transactionID)
}
