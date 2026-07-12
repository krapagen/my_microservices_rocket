package v1

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func New(partService PartService) *api {
	return &api{
		partService: partService,
	}
}

type Converter interface {
	ToGetInput(rawUUID string) (uuid.UUID, error)
	PartToProto(part model.Part) *inventoryv1.Part
	PartsToProto(parts []model.Part) []*inventoryv1.Part
	PartTypeToProtoPartType(t model.PartType) inventoryv1.PartType
	ProtoPartTypeToPartType(t inventoryv1.PartType) model.PartType
}
