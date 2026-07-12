package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (r *RepoSuite) TestGet_Success() {
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

	err := r.repo.Create(r.ctx, testOrder)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)

	r.Equal(testOrder.UUID, retrieved.UUID)
	r.Equal(testOrder.Status, retrieved.Status)
	r.Equal(len(testOrder.Items), len(retrieved.Items))
	r.Equal(testOrder.Items[0].PartUUID, retrieved.Items[0].PartUUID)
	r.Equal(testOrder.Items[0].PartType, retrieved.Items[0].PartType)
	r.Equal(testOrder.Items[0].Price, retrieved.Items[0].Price)
}

func (r *RepoSuite) TestGet_NotFound() {
	nonexistentOrderID := uuid.New()

	_, err := r.repo.Get(r.ctx, nonexistentOrderID)
	r.Error(err)
	r.True(errors.Is(err, errs.ErrOrderNotFound))
}

func (r *RepoSuite) TestGet_MultipleItems() {
	orderID := uuid.New()
	partID1 := uuid.New()
	partID2 := uuid.New()
	partID3 := uuid.New()

	testOrder := model.Order{
		UUID: orderID,
		Items: []model.OrderItem{
			{
				PartUUID: partID1,
				PartType: model.PartTypeHull,
				Price:    1000,
			},
			{
				PartUUID: partID2,
				PartType: model.PartTypeEngine,
				Price:    500,
			},
			{
				PartUUID: partID3,
				PartType: model.PartTypeShield,
				Price:    300,
			},
		},
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	err := r.repo.Create(r.ctx, testOrder)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)

	r.Equal(3, len(retrieved.Items))
	r.Equal(testOrder.Items[0].PartUUID, retrieved.Items[0].PartUUID)
	r.Equal(testOrder.Items[1].PartUUID, retrieved.Items[1].PartUUID)
	r.Equal(testOrder.Items[2].PartUUID, retrieved.Items[2].PartUUID)
}

func (r *RepoSuite) TestGet_WithPaymentInfo() {
	orderID := uuid.New()
	partID := uuid.New()
	transactionID := uuid.New()
	paymentMethod := model.PaymentMethodCard

	testOrder := model.Order{
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

	err := r.repo.Create(r.ctx, testOrder)
	r.NoError(err)

	retrieved, err := r.repo.Get(r.ctx, orderID)
	r.NoError(err)

	r.Equal(testOrder.Status, retrieved.Status)
	r.NotNil(retrieved.TransactionUUID)
	r.Equal(transactionID, *retrieved.TransactionUUID)
	r.NotNil(retrieved.PaymentMethod)
	r.Equal(paymentMethod, *retrieved.PaymentMethod)
}
