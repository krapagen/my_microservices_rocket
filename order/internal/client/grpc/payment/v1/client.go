package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

// New создаёт обёртку над gRPC клиентом PaymentService.
func New(c paymentv1.PaymentServiceClient) *client {
	return &client{paymentClient: c}
}

func (c client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	op := "order/internal/client/grpc/payment/v1/Payorder"
	log := slog.With("op", op)
	payment, err := converter.PaymentToProto(method)
	if err != nil {
		log.ErrorContext(ctx, "Ошибка конвертации model в proto")
		return uuid.Nil, err
	}
	log.InfoContext(ctx, "Отправка запроса на оплату заказа", "orderUUID", orderUUID.String(), "paymentMethod", payment.String())
	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: payment,
	})
	if err != nil {
		// Обрабатываем gRPC ошибки
		switch status.Code(err) {
		case codes.NotFound:
			log.ErrorContext(ctx, "Заказ не найден", "orderUUID", orderUUID.String())
			return uuid.Nil, errs.ErrOrderNotFound
		case codes.FailedPrecondition:
			log.ErrorContext(ctx, "Заказ уже оплачен", "orderUUID", orderUUID.String())
			return uuid.Nil, errs.ErrOrderAlreadyPaid
		case codes.Aborted:
			log.ErrorContext(ctx, "Заказ отменён", "orderUUID", orderUUID.String())
			return uuid.Nil, errs.ErrOrderCancelled
		default:
			return uuid.Nil, fmt.Errorf("ошибка оплаты: %w", err)
		}
	}
	log.InfoContext(ctx, "Заказ успешно оплачен", "orderUUID", orderUUID.String(), "transactionUUID", resp.GetTransactionUuid())
	// Парсим UUID транзакции из ответа
	transactionUUID, err := converter.GetTransaction(resp)
	if err != nil {
		log.ErrorContext(ctx, "Ошибка конвертации transactionUUID из proto в model")
		return uuid.Nil, err
	}
	log.InfoContext(ctx, "UUID транзакции успешно получен", "transactionUUID", transactionUUID.String())
	return transactionUUID, nil
}
