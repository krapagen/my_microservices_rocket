package converter

import (
	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/order/internal/repository/record"
)

func OrderToModel(order record.Order, orderItems []record.OrderItem) model.Order {
	return model.Order{
		UUID:            order.UUID,
		Items:           ItemsRecordToModel(orderItems),
		TransactionUUID: TransactionRecordToModel(order.TransactionUUID),
		PaymentMethod:   PaymentMethodRecordToModel(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
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
		PartUUID: item.PartUUID,
		PartType: model.PartType(item.PartType),
		Price:    item.Price,
	}
}

func PaymentMethodRecordToModel(method *string) *model.PaymentMethod {
	if method == nil {
		return nil
	}
	return new(model.PaymentMethod(*method))
}

func TransactionRecordToModel(transactionUUID *uuid.UUID) *uuid.UUID {
	if transactionUUID == nil {
		return nil
	}
	return new(*transactionUUID)
}

func OrderToRecord(order model.Order) (record.Order, []record.OrderItem) {
	return record.Order{
		UUID:            order.UUID,
		TransactionUUID: TransactionModelToRecord(order.TransactionUUID),
		PaymentMethod:   PaymentMethodModelToRecord(order.PaymentMethod),
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
	}, ItemsModelToRecord(order.Items, order.UUID)
}

func ItemsModelToRecord(items []model.OrderItem, orderUUID uuid.UUID) []record.OrderItem {
	recordItems := make([]record.OrderItem, 0, len(items))
	for _, item := range items {
		recordItems = append(recordItems, ItemModelToRecord(item, orderUUID))
	}
	return recordItems
}

func ItemModelToRecord(item model.OrderItem, orderUUID uuid.UUID) record.OrderItem {
	return record.OrderItem{
		OrderUUID: orderUUID,
		PartUUID:  item.PartUUID,
		PartType:  string(item.PartType),
		Price:     item.Price,
	}
}

func PaymentMethodModelToRecord(method *model.PaymentMethod) *string {
	if method == nil {
		return nil
	}
	return new(string(*method))
}

func TransactionModelToRecord(transactionUUID *uuid.UUID) *uuid.UUID {
	if transactionUUID == nil {
		return nil
	}

	return new(*transactionUUID)
}
