package v1

import (
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
