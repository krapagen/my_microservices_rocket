package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/krapagen/my_microservices_rocket/payment/internal/api/payment/v1"
	"github.com/krapagen/my_microservices_rocket/payment/internal/api/payment/v1/mocks"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

type APISuite struct {
	suite.Suite
	ctx            context.Context
	paymentService *mocks.PaymentService
	api            paymentv1.PaymentServiceServer
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.paymentService = mocks.NewPaymentService(s.T())
	s.api = v1.NewAPI(s.paymentService)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
