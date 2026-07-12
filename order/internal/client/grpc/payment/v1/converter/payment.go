package converter

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

func PaymentToProto(payment model.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch payment {
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, errs.ErrInvalidPaymentMethod
	}
}

// Парсим UUID транзакции из ответа
func GetTransaction(resp *paymentv1.PayOrderResponse) (uuid.UUID, error) {
	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("неверный формат UUID транзакции: %w", errs.ErrInvalidUUID)
	}
	return transactionUUID, nil
}
