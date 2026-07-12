package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	inventoryv1 "github.com/krapagen/my_microservices_rocket/inventory/internal/api/inventory/v1"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/part"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/part/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx            context.Context
	partRepository *mocks.PartRepository
	service        inventoryv1.PartService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.partRepository = mocks.NewPartRepository(s.T())
	s.service = part.New(s.partRepository)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
