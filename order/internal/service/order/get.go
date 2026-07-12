package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	op := "order/internal/service/order/Get"
	log := slog.With("op", op)
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		log.ErrorContext(ctx, "не удалось получить заказ")
		if errors.Is(err, errs.ErrOrderNotFound) {
			log.ErrorContext(ctx, "Заказ не найден", "orderUUID", orderUUID)
			return model.Order{}, errs.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}
	log.InfoContext(ctx, "Заказ успешно получен", "orderUUID", orderUUID)
	return order, nil
}
