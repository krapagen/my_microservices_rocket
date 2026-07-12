package converter

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/service/input"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

func OrderModelToDTO(order model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
		Status:     orderv1.OrderStatus(order.Status),
		CreatedAt:  order.CreatedAt,
	}

	// Извлекаем UUID деталей из Items
	hullUUID, engineUUID, shieldUUID, weaponUUID := extractPartUUIDs(order.Items)

	dto.HullUUID = hullUUID
	dto.EngineUUID = engineUUID

	// Опциональные поля
	if shieldUUID != uuid.Nil {
		dto.ShieldUUID = orderv1.OptNilUUID{
			Value: shieldUUID,
			Set:   true,
			Null:  false,
		}
	} else {
		dto.ShieldUUID = orderv1.OptNilUUID{
			Null: true,
		}
	}

	if weaponUUID != uuid.Nil {
		dto.WeaponUUID = orderv1.OptNilUUID{
			Value: weaponUUID,
			Set:   true,
			Null:  false,
		}
	} else {
		dto.WeaponUUID = orderv1.OptNilUUID{
			Null: true,
		}
	}

	// Опциональные поля после оплаты
	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderv1.OptNilUUID{
			Value: *order.TransactionUUID,
			Set:   true,
			Null:  false,
		}
	} else {
		dto.TransactionUUID = orderv1.OptNilUUID{
			Null: true,
		}
	}

	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderv1.OptNilPaymentMethod{
			Value: orderv1.PaymentMethod(*order.PaymentMethod),
			Set:   true,
			Null:  false,
		}
	} else {
		dto.PaymentMethod = orderv1.OptNilPaymentMethod{
			Null: true,
		}
	}

	return dto
}

func extractPartUUIDs(items []model.OrderItem) (hullUUID, engineUUID, shieldUUID, weaponUUID uuid.UUID) {
	for _, item := range items {
		switch item.PartType {
		case model.PartTypeHull:
			hullUUID = item.PartUUID
		case model.PartTypeEngine:
			engineUUID = item.PartUUID
		case model.PartTypeShield:
			shieldUUID = item.PartUUID
		case model.PartTypeWeapon:
			weaponUUID = item.PartUUID
		}
	}

	return hullUUID, engineUUID, shieldUUID, weaponUUID
}

func CreateOrderRequestToInput(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	orderInput := input.CreateOrderInput{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
	}

	if req.ShieldUUID.Set && !req.ShieldUUID.Null {
		orderInput.ShieldUUID = &req.ShieldUUID.Value
	}

	if req.WeaponUUID.Set && !req.WeaponUUID.Null {
		orderInput.WeaponUUID = &req.WeaponUUID.Value
	}

	return orderInput
}

// PayOrderRequestToModel converts PayOrder request to service parameters
func PayOrderRequestToModel(req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
} {
	return struct {
		OrderUUID     uuid.UUID
		PaymentMethod model.PaymentMethod
	}{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
	}
}
