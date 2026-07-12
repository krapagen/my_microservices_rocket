package tests

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (s *ServiceSuite) TestGet_Success() {
	orderID := uuid.New()
	partID := uuid.New()
	createdAt := time.Now()

	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: createdAt,
	}
	s.orderRepository.EXPECT().Get(s.ctx, orderID).Return(testOrder, nil)

	result, err := s.service.Get(s.ctx, orderID)
	s.NoError(err)
	s.Equal(testOrder.UUID, result.UUID)
	s.Equal(testOrder.Status, result.Status)
	s.Equal(testOrder.Items, result.Items)
	s.Equal(testOrder.CreatedAt, result.CreatedAt)

}

func (s *ServiceSuite) TestGet_NotFound() {

	orderUUID := uuid.New()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{}, errs.ErrOrderNotFound)

	result, err := s.service.Get(s.ctx, orderUUID)

	s.Error(err)
	s.ErrorIs(err, errs.ErrOrderNotFound)
	s.Equal(model.Order{}, result)
}

func (s *ServiceSuite) TestGet_RepoError() {
	var (
		orderUUID = uuid.New()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{}, repoErr)

	result, err := s.service.Get(s.ctx, orderUUID)

	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Equal(model.Order{}, result)
}
