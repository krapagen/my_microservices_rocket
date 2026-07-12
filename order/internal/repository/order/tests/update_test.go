package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (r *RepoSuite) TestUpdate_Success() {
	orderID := uuid.New()
	partID := uuid.New()
	createdAt := time.Now()

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
		CreatedAt: createdAt,
	}

	err := r.repo.Create(r.ctx, order)
	r.NoError(err)

	updatedOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID,
				PartType: model.PartTypeEngine,
				Price:    1500,
			},
		},
		Status:    model.OrderStatusPaid,
		CreatedAt: createdAt,
	}

	err = r.repo.Update(r.ctx, updatedOrder)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)
	r.Equal(updatedOrder.UUID, retrieved.UUID)
	r.Equal(updatedOrder.Status, retrieved.Status)
	r.Equal(len(updatedOrder.Items), len(retrieved.Items))
	r.Equal(updatedOrder.Items[0].PartUUID, retrieved.Items[0].PartUUID)
	r.Equal(updatedOrder.Items[0].PartType, retrieved.Items[0].PartType)
	r.Equal(updatedOrder.Items[0].Price, retrieved.Items[0].Price)
}

func (r *RepoSuite) TestUpdate_NotFound() {
	nonexistentOrder := model.Order{
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

	err := r.repo.Update(r.ctx, nonexistentOrder)
	r.Error(err)
	r.True(errors.Is(err, errs.ErrOrderNotFound))
}

func (r *RepoSuite) TestUpdate_WithPaymentInfo() {
	orderID := uuid.New()
	partID := uuid.New()
	createdAt := time.Now()

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
		CreatedAt: createdAt,
	}

	err := r.repo.Create(r.ctx, order)
	r.NoError(err)

	transactionID := uuid.New()
	paymentMethod := model.PaymentMethodCard

	updatedOrder := model.Order{
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
		CreatedAt:       createdAt,
	}

	err = r.repo.Update(r.ctx, updatedOrder)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)
	r.Equal(updatedOrder.Status, retrieved.Status)
	r.NotNil(retrieved.TransactionUUID)
	r.Equal(transactionID, *retrieved.TransactionUUID)
	r.NotNil(retrieved.PaymentMethod)
	r.Equal(paymentMethod, *retrieved.PaymentMethod)
}
