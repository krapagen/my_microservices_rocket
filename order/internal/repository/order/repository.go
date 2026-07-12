package order

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	repoConverter "github.com/krapagen/my_microservices_rocket/order/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/order/internal/repository/record"
)

type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

// New создаёт новый in-memory репозиторий заказов.
func New() *repository {
	return &repository{
		orders: make(map[uuid.UUID]record.Order),
	}
}

func (r *repository) Create(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/Create"
	log := slog.With("op", op)
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.UUID] = repoConverter.OrderToRepoModel(order)
	log.InfoContext(ctx, "Заказ Создан", "order.UUID", order.UUID)
	return nil
}

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	op := "order/internal/repository/order/Get"
	log := slog.With("op", op)
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.orders[orderUUID]
	if !ok {
		log.ErrorContext(ctx, "Заказ не найден", "orderUUID", orderUUID)
		return model.Order{}, errs.ErrOrderNotFound
	}
	log.InfoContext(ctx, "Заказ найден", "orderUUID", orderUUID)
	return repoConverter.OrderToModel(rec), nil
}

func (r *repository) Update(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/Update"
	log := slog.With("op", op)
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.UUID]; !ok {
		log.ErrorContext(ctx, "Заказ не найден для изменения", "orderUUID", order.UUID)
		return errs.ErrOrderNotFound
	}
	r.orders[order.UUID] = repoConverter.OrderToRepoModel(order)
	log.InfoContext(ctx, "Заказ изменен", "orderUUID", order.UUID)
	return nil
}
