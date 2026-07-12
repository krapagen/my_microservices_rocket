package converter

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/record"
)

func PartRecordToModel(part record.Part) model.Part {
	return model.Part{
		UUID:          uuid.MustParse(part.UUID),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      RecordToModelType[part.PartType],
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}

var RecordToModelType = map[record.PartType]model.PartType{
	record.PartType_PART_TYPE_UNSPECIFIED: model.PartTypeUnspecified,
	record.PartType_PART_TYPE_HULL:        model.PartTypeHull,
	record.PartType_PART_TYPE_ENGINE:      model.PartTypeEngine,
	record.PartType_PART_TYPE_SHIELD:      model.PartTypeShield,
	record.PartType_PART_TYPE_WEAPON:      model.PartTypeWeapon,
}
