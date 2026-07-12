package v1

import orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"

type api struct {
	orderv1.UnimplementedHandler

	orderService OrderService
}

func NewAPI(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
