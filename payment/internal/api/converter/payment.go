package converter

import (
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/payment/internal/errors"
	"github.com/krapagen/my_microservices_rocket/payment/internal/model"
	"github.com/krapagen/my_microservices_rocket/payment/internal/service/input"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

func PayOrderRequestProtoToInput(req *paymentv1.PayOrderRequest) (input.PayOrderInput, error) {
	parsedUuid, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return input.PayOrderInput{}, errs.ErrInvalidOrderUUID
	}
	paymentMethod, err := ToModelPayment(req.GetPaymentMethod())
	if err != nil {
		return input.PayOrderInput{}, err
	}
	return input.PayOrderInput{
		OrderUUID:     parsedUuid,
		PaymentMethod: paymentMethod,
	}, nil
}

func ToModelPayment(method paymentv1.PaymentMethod) (model.PaymentMethod, error) {
	switch method {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnspecified, errs.ErrInvalidPaymentMethod
	}
}
