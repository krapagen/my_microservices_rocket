package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/krapagen/my_microservices_rocket/payment/internal/service/input"
	"github.com/krapagen/my_microservices_rocket/payment/internal/service/payment"
)

type PaymentService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
}

type ServiceSuite struct {
	suite.Suite
	ctx     context.Context
	service PaymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = payment.New()
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
