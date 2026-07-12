package input

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}
