package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	orderv1 "github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/order"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/order/mocks"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/auth"
)

type ServiceSuite struct {
	suite.Suite
	ctx                  context.Context
	orderRepository      *mocks.OrderRepository
	orderPaymentClient   *mocks.PaymentClient
	orderInventoryClient *mocks.InventoryClient
	txManager            *mocks.TxManager
	service              orderv1.OrderService
}

type noopOrderPaidProducer struct{}

func (noopOrderPaidProducer) ProduceOrderPaid(_ context.Context, _ model.OrderPaid) error {
	return nil
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = auth.WithUserUUID(context.Background(), uuid.New())
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.orderPaymentClient = mocks.NewPaymentClient(s.T())
	s.orderInventoryClient = mocks.NewInventoryClient(s.T())
	s.txManager = mocks.NewTxManager(s.T())
	s.txManager.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	}).Maybe()
	s.service = order.New(
		s.orderRepository,
		s.orderInventoryClient,
		s.orderPaymentClient,
		noopOrderPaidProducer{},
		s.txManager,
	)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
