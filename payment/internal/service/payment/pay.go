package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/payment/internal/errors"
	"github.com/krapagen/my_microservices_rocket/payment/internal/service/input"
)

func (s *service) Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	// 1. Валидация payment_method
	//    (order_uuid уже распарсен в api/converter и в сервис приходит как uuid.UUID)
	op := "payment/internal/service/payment/Pay"
	log := slog.With("op", op)
	if !in.PaymentMethod.IsValid() {
		log.ErrorContext(ctx, "неверный payment_method", "payment_method", in.PaymentMethod)
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}
	log.InfoContext(ctx, "начало оплаты", "order_uuid", in.OrderUUID, "payment_method", in.PaymentMethod)
	// 2. Генерация transaction_uuid
	transactionUUID := uuid.New()

	// 3. Логирование
	log.InfoContext(ctx, "оплата выполнена",
		"order_uuid", in.OrderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", in.PaymentMethod)

	return transactionUUID, nil
}
