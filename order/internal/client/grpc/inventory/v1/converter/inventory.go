package converter

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func PartTypeFromProto(part inventoryv1.PartType) model.PartType {
	switch part {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	}
	return ""
}

func PartsToModel(parts []*inventoryv1.Part) []model.Part {
	res := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		res = append(res, model.Part{
			UUID:          uuid.MustParse(p.GetUuid()),
			Name:          p.GetName(),
			PartType:      PartTypeFromProto(p.GetPartType()),
			Price:         p.GetPrice(),
			StockQuantity: p.GetStockQuantity(),
		})
	}
	return res
}
