package part

import "github.com/krapagen/my_microservices_rocket/inventory/internal/service/domain"

type service struct {
	partRepository       PartRepository
	compatibilityChecker CompatibilityChecker
}

func New(partRepository PartRepository) *service {
	return &service{
		partRepository:       partRepository,
		compatibilityChecker: domain.NewCompatibilityChecker(),
	}
}
