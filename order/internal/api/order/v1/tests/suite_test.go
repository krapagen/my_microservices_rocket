package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	"github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1/mocks"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

type ServiceSuite struct {
	suite.Suite
	ctx          context.Context
	orderService *mocks.OrderService
	api          orderv1.Handler
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderService = mocks.NewOrderService(s.T())
	s.api = v1.NewAPI(s.orderService)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
