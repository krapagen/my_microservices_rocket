package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

type converter struct{}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToGetInput(rawUUID string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return parsed, nil
}

func (c *converter) ToGetInputs(rawUUIDs []string) ([]uuid.UUID, error) {
	if len(rawUUIDs) == 0 {
		return make([]uuid.UUID, 0), nil
	}
	curUuids := make([]uuid.UUID, 0, len(rawUUIDs))
	for _, rawUUID := range rawUUIDs {
		curUuid, err := c.ToGetInput(rawUUID)
		if err != nil {
			return nil, err
		}
		curUuids = append(curUuids, curUuid)
	}
	return curUuids, nil
}

func (c *converter) PartToDto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID().String(),
		Name:          part.Name(),
		Description:   part.Description(),
		Price:         part.Price(),
		PartType:      c.PartTypeToDtoType(part.PartType()),
		StockQuantity: int64(part.StockQuantity()),
		CreatedAt:     timestamppb.New(part.CreatedAt()),
	}
}

func (c *converter) PartsToDto(parts []model.Part) []*inventoryv1.Part {
	res := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		res = append(res, c.PartToDto(p))
	}
	return res
}

func (c *converter) PartTypeToDtoType(t model.PartType) inventoryv1.PartType {
	switch t {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}

func (c *converter) DtoTypeToPartType(t inventoryv1.PartType) model.PartType {
	switch t {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}
