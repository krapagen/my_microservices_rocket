package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
	"github.com/stretchr/testify/suite"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/part"
)

type Repository interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
}

type RepoSuite struct {
	suite.Suite
	ctx  context.Context
	repo Repository
}

func (r *RepoSuite) SetupTest() {
	r.ctx = context.Background()
	r.repo = part.NewRepository()
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
