package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/repository/order"
)

type Repository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type RepoSuite struct {
	suite.Suite
	ctx  context.Context
	repo Repository
}

func (r *RepoSuite) SetupTest() {
	r.ctx = context.Background()
	r.repo = order.New()
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
