package assembly

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
)

// Assembler определяет контракт бизнес-логики сборки корабля.
type Assembler interface {
	Assemble(ctx context.Context, in model.OrderPaid) (model.ShipAssembled, error)
}