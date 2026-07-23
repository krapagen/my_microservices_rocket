package order

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
)

type service struct {
	orderRepo       OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       TxManager
}

// New создаёт новый сервис заказов.
func New(
	orderRepo OrderRepository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	txManager TxManager,
) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		txManager:       txManager,
	}
}

func partUUIDsFromOrderItems(items []model.OrderItem) []uuid.UUID {
	uuids := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		uuids = append(uuids, item.PartUUID)
	}
	return uuids
}
