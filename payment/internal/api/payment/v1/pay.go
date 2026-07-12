package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/payment/internal/api/converter"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	op := "payment/internal/api/v1//PayOrder"
	log := slog.With("op", op)
	modelReq, err := converter.PayOrderRequestProtoToInput(req)
	if err != nil {
		log.ErrorContext(ctx, "не удалось конвертировать запрос ошибка uuid или paymentMethod", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "входные данные верны uuid и paymentMethod", "request", modelReq)
	transactionUUID, err := a.paymentService.Pay(ctx, modelReq)
	if err != nil {
		log.ErrorContext(ctx, "Не удалось оплатить заказ", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "оплата заказа прошла успешно", "transactionUUID", transactionUUID)
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
