package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	v1 "github.com/krapagen/my_microservices_rocket/inventory/internal/api/inventory/v1"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/api/inventory/v1/mocks"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

type APISuite struct {
	suite.Suite
	ctx         context.Context
	partService *mocks.PartService
	api         inventoryv1.InventoryServiceServer
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.partService = mocks.NewPartService(s.T())
	s.api = v1.New(s.partService)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
