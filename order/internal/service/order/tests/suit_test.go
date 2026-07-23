package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/input"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/order"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/order/mocks"
)

type OrderService interface {
	Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error)
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
}

type ServiceSuite struct {
	suite.Suite
	ctx                  context.Context
	orderRepository      *mocks.OrderRepository
	orderPaymentClient   *mocks.PaymentClient
	orderInventoryClient *mocks.InventoryClient
	txManager            *mocks.TxManager
	service              OrderService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.orderPaymentClient = mocks.NewPaymentClient(s.T())
	s.orderInventoryClient = mocks.NewInventoryClient(s.T())
	s.txManager = mocks.NewTxManager(s.T())
	s.txManager.EXPECT().Do(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	}).Maybe()
	s.service = order.New(s.orderRepository, s.orderInventoryClient, s.orderPaymentClient, s.txManager)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
