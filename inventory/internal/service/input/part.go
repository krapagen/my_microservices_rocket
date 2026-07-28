package input

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

type PartFilter struct {
	// UUIDs — если не пустой, возвращаются только эти детали (приоритет)
	UUIDs []uuid.UUID
	// PartType — фильтр по типу (игнорируется если UUIDs заполнен)
	PartType model.PartType
}

type CommitFilter struct {
	// UUIDs - эти детали будут забронированы, затем сняты с резерва и списаны со склада
	UUIDs []uuid.UUID
}
