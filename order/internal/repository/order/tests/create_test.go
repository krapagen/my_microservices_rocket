package tests

import (
	"time"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (r *RepoSuite) TestCreate_Success() {
	orderID := uuid.New()
	partID := uuid.New()
	now := time.Now()

	order := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: now,
	}

	err := r.repo.Create(r.ctx, order)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)
	r.Equal(order.UUID, retrieved.UUID)
	r.Equal(order.Status, retrieved.Status)
}

func (r *RepoSuite) TestCreate_MultipleOrders() {
	order1 := model.Order{
		UUID: uuid.New(),
		Items: []model.OrderItem{
			{
				PartUUID: uuid.New(),
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	order2 := model.Order{
		UUID: uuid.New(),
		Items: []model.OrderItem{
			{
				PartUUID: uuid.New(),
				PartType: model.PartTypeEngine,
				Price:    500,
			},
		},
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	err := r.repo.Create(r.ctx, order1)
	r.NoError(err)

	err = r.repo.Create(r.ctx, order2)
	r.NoError(err)

	retrieved1, err := r.repo.Get(r.ctx, order1.UUID)
	r.NoError(err)
	r.Equal(order1.UUID, retrieved1.UUID)

	retrieved2, err := r.repo.Get(r.ctx, order2.UUID)
	r.NoError(err)
	r.Equal(order2.UUID, retrieved2.UUID)
}

func (r *RepoSuite) TestCreate_WithPaymentInfo() {
	orderID := uuid.New()
	partID := uuid.New()
	transactionID := uuid.New()
	paymentMethod := model.PaymentMethodCard

	order := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
		},
		TransactionUUID: &transactionID,
		PaymentMethod:   &paymentMethod,
		Status:          model.OrderStatusPaid,
		CreatedAt:       time.Now(),
	}

	err := r.repo.Create(r.ctx, order)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)
	r.NotNil(retrieved.TransactionUUID)
	r.Equal(transactionID, *retrieved.TransactionUUID)
	r.NotNil(retrieved.PaymentMethod)
	r.Equal(paymentMethod, *retrieved.PaymentMethod)
}
