package converter

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/repository/record"
)

func OrderToModel(order record.Order) model.Order {
	return model.Order{
		UUID:            uuid.MustParse(order.OrderUUID),
		Items:           ItemsRecordToModel(order.Items),
		TransactionUUID: TransactionRecordToModel(order.TransactionUUID),
		PaymentMethod:   PaymentMethodRecordToModel(order.PaymentMethod),
		Status:          StatusRecordToModel(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func ItemsRecordToModel(items []record.OrderItem) []model.OrderItem {
	modelItems := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, ItemRecordToModel(item))
	}
	return modelItems
}

func ItemRecordToModel(item record.OrderItem) model.OrderItem {
	return model.OrderItem{
		PartUUID: uuid.MustParse(item.PartUUID),
		PartType: PartTypeRecordToModel(item.PartType),
		Price:    item.Price,
	}
}

func PartTypeRecordToModel(part record.PartType) model.PartType {
	switch part {
	case record.PartTypeHull:
		return model.PartTypeHull
	case record.PartTypeEngine:
		return model.PartTypeEngine
	case record.PartTypeShield:
		return model.PartTypeShield
	case record.PartTypeWeapon:
		return model.PartTypeWeapon
	default:
		return ""
	}
}

func StatusRecordToModel(status record.OrderStatus) model.OrderStatus {
	switch status {
	case record.OrderStatusPendingPayment:
		return model.OrderStatusPendingPayment
	case record.OrderStatusPaid:
		return model.OrderStatusPaid
	case record.OrderStatusCancelled:
		return model.OrderStatusCancelled
	default:
		return ""
	}
}

func PaymentMethodRecordToModel(method *record.PaymentMethod) *model.PaymentMethod {
	if method == nil {
		return nil
	}
	return new(model.PaymentMethod(*method))
}

func TransactionRecordToModel(transactionUUID *string) *uuid.UUID {
	if transactionUUID == nil {
		return nil
	}
	return new(uuid.MustParse(*transactionUUID))
}

func OrderToRepoModel(order model.Order) record.Order {
	return record.Order{
		OrderUUID:       order.UUID.String(),
		Items:           ItemsModelToRecord(order.Items),
		TransactionUUID: TransactionModelToRecord(order.TransactionUUID),
		PaymentMethod:   PaymentMethodModelToRecord(order.PaymentMethod),
		Status:          StatusModelToRecord(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func ItemsModelToRecord(items []model.OrderItem) []record.OrderItem {
	recordItems := make([]record.OrderItem, 0, len(items))
	for _, item := range items {
		recordItems = append(recordItems, ItemModelToRecord(item))
	}
	return recordItems
}

func ItemModelToRecord(item model.OrderItem) record.OrderItem {
	return record.OrderItem{
		PartUUID: item.PartUUID.String(),
		PartType: PartTypeModelToRecord(item.PartType),
		Price:    item.Price,
	}
}

func PartTypeModelToRecord(part model.PartType) record.PartType {
	switch part {
	case model.PartTypeHull:
		return record.PartTypeHull
	case model.PartTypeEngine:
		return record.PartTypeEngine
	case model.PartTypeShield:
		return record.PartTypeShield
	case model.PartTypeWeapon:
		return record.PartTypeWeapon
	default:
		return ""
	}
}

func StatusModelToRecord(status model.OrderStatus) record.OrderStatus {
	switch status {
	case model.OrderStatusPendingPayment:
		return record.OrderStatusPendingPayment
	case model.OrderStatusPaid:
		return record.OrderStatusPaid
	case model.OrderStatusCancelled:
		return record.OrderStatusCancelled
	default:
		return ""
	}
}

func PaymentMethodModelToRecord(method *model.PaymentMethod) *record.PaymentMethod {
	if method == nil {
		return nil
	}
	return new(record.PaymentMethod(*method))
}

func TransactionModelToRecord(transactionUUID *uuid.UUID) *string {
	if transactionUUID == nil {
		return nil
	}

	return new(transactionUUID.String())
}
