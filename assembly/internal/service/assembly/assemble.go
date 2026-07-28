package assembly

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
)

func (s *service) Assemble(ctx context.Context, in model.OrderPaid) (model.ShipAssembled, error) {
	timer := time.NewTimer(s.buildTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return model.ShipAssembled{}, ctx.Err()
	case <-timer.C:
	}

	now := time.Now().UTC()

	return model.ShipAssembled{
		EventUUID:   uuid.New(),
		OrderUUID:   in.OrderUUID,
		UserUUID:    in.UserUUID,
		BuildTime:   s.buildTime,
		AssembledAt: now,
	}, nil
}
