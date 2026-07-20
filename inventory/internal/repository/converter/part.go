package converter

import (
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/record"
)

func PartRecordToModel(part record.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}

func PartsRecordToModel(parts []record.Part) []model.Part {
	modelParts := make([]model.Part, len(parts))
	for i, part := range parts {
		modelParts[i] = PartRecordToModel(part)
	}
	return modelParts
}
