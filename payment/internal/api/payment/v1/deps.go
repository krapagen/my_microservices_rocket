package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/payment/internal/service/input"
)

type PaymentService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
}
