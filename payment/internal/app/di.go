package app

import (
	"context"

	paymentApi "github.com/krapagen/my_microservices_rocket/payment/internal/api/payment/v1"
	paymentService "github.com/krapagen/my_microservices_rocket/payment/internal/service/payment"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	// Сервисы
	paymentService paymentApi.PaymentService

	// API-обработчики
	paymentV1Handler paymentv1.PaymentServiceServer
}

func (d *diContainer) PaymentService(ctx context.Context) paymentApi.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.New()
	}
	return d.paymentService
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentV1Handler == nil {
		d.paymentV1Handler = paymentApi.NewAPI(d.PaymentService(ctx))
	}
	return d.paymentV1Handler
}
