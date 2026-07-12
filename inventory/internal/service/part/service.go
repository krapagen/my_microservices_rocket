package part

type service struct {
	partRepository PartRepository
}

func New(partRepository PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
