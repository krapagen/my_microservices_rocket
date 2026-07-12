package order

type service struct {
	orderRepo       OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

// New создаёт новый сервис заказов.
func New(
	orderRepo OrderRepository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
