package part

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/domain"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
type service struct {
	partRepository       PartRepository
	compatibilityChecker CompatibilityChecker
	txManager            TxManager
}

func New(partRepository PartRepository, txManager TxManager) *service {
	return &service{
		partRepository:       partRepository,
		compatibilityChecker: domain.NewCompatibilityChecker(),
		txManager:            txManager,
	}
}
