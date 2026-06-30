package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты заказа
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// orderStore — хранилище заказов (in-memory)
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// handler реализует интерфейс orderv1.Handler, сгенерированный ogen
type handler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *orderStore
}

// NewHandler создаёт новый обработчик заказов
func NewHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *orderStore,
) *handler {
	return &handler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика
func SetupServer(h *handler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации)
// GET /api/v1/orders/{order_uuid}.
func (h *handler) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	log := slog.With("op", "order/pkg/handler/GetOrder")

	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		log.ErrorContext(ctx, "заказ не найден", "order_uuid", params.OrderUUID)
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	log.InfoContext(ctx, "заказ найден", "order_uuid", params.OrderUUID, "status", order.Status)

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CreateOrder реализует операцию createOrder
// POST /api/v1/orders
func (h *handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	log := slog.With("op", "order/pkg/handler/CreateOrder")

	// 1. Валидация: hull_uuid и engine_uuid обязательны
	hullUUID := req.GetHullUUID()
	engineUUID := req.GetEngineUUID()

	if hullUUID == uuid.Nil || engineUUID == uuid.Nil {
		log.ErrorContext(ctx, "hull_uuid и engine_uuid обязательны")
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid и engine_uuid обязательны",
		}, nil
	}

	uuids := []string{
		hullUUID.String(),
		engineUUID.String(),
	}
	var shieldUUIDPtr *uuid.UUID
	if shieldUUID, ok := req.ShieldUUID.Get(); ok {
		uuids = append(uuids, shieldUUID.String())
		shieldUUIDPtr = &shieldUUID
	}
	var weaponUUIDPtr *uuid.UUID
	if weaponUUID, ok := req.WeaponUUID.Get(); ok {
		uuids = append(uuids, weaponUUID.String())
		weaponUUIDPtr = &weaponUUID
	}

	listReq := &inventoryv1.ListPartsRequest{Uuids: uuids}

	// 2. Получить детали через InventoryService.ListParts
	listResp, err := h.inventoryClient.ListParts(ctx, listReq)
	if err != nil {
		log.ErrorContext(ctx, "failed to get parts from inventory", "error", err)
		// Проверяем статус gRPC ошибки
		if status.Code(err) == codes.NotFound {
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "part not found",
			}, nil
		}
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get parts from inventory",
		}, nil
	}

	// 3. Проверить stock_quantity > 0
	partsMap := make(map[string]*inventoryv1.Part)
	for _, part := range listResp.Parts {
		if part.StockQuantity <= 0 {
			log.ErrorContext(ctx, "part is out of stock", "part_uuid", part.Uuid)
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("part %s is out of stock", part.Uuid),
			}, nil
		}
		partsMap[part.Uuid] = part
	}

	// Validate required parts exist
	if _, ok := partsMap[hullUUID.String()]; !ok {
		log.ErrorContext(ctx, "hull not found", "hull_uuid", hullUUID.String())
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "hull not found",
		}, nil
	}
	if _, ok := partsMap[engineUUID.String()]; !ok {
		log.ErrorContext(ctx, "engine not found", "engine_uuid", engineUUID.String())
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "engine not found",
		}, nil
	}

	// 4. Вычислить total_price
	var totalPrice int64
	for _, part := range listResp.Parts {
		totalPrice += part.Price
	}

	// 5. Сгенерировать order_uuid (UUID v4)
	orderUuid := uuid.New()

	// 6. Создать заказ со статусом PENDING_PAYMENT
	orderCurrent := &Order{
		OrderUUID:       orderUuid,
		HullUUID:        hullUUID,
		EngineUUID:      engineUUID,
		ShieldUUID:      shieldUUIDPtr,
		WeaponUUID:      weaponUUIDPtr,
		TotalPrice:      totalPrice,
		TransactionUUID: nil, // пока нет транзакции
		PaymentMethod:   nil, // пока нет способа оплаты
		Status:          OrderStatusPendingPayment,
		CreatedAt:       time.Now(),
	}
	// 7. Сохранить в store
	h.store.mu.Lock()
	h.store.orders[orderUuid] = *orderCurrent
	h.store.mu.Unlock()

	log.InfoContext(ctx, "заказ создан", "order_uuid", orderUuid, "total_price", totalPrice)

	// 8. Вернуть order_uuid и total_price
	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay
func (h *handler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	log := slog.With("op", "order/pkg/handler/PayOrder")

	// 1. Найти заказ в store

	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// Если не найден — вернуть 404
	if !ok {
		log.ErrorContext(ctx, "заказ не найден", "order_uuid", params.OrderUUID)
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("заказ %s не найден", params.OrderUUID),
		}, nil
	}

	// 2. Проверить статус == PENDING_PAYMENT

	if order.Status != OrderStatusPendingPayment {
		log.ErrorContext(ctx, "заказ не в статусе PENDING_PAYMENT", "order_uuid", params.OrderUUID, "status", order.Status)
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("заказ %s не в статусе PENDING_PAYMENT, его статус %s", params.OrderUUID, order.Status),
		}, nil
	}

	// 3. Вызвать h.paymentClient.PayOrder для обработки платежа

	var protoPaymentMethod paymentv1.PaymentMethod
	switch req.PaymentMethod {
	case orderv1.PaymentMethodCARD:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		protoPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		log.ErrorContext(ctx, "неизвестный способ оплаты", "payment_method", req.PaymentMethod)
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("неизвестный способ оплаты %s", req.PaymentMethod),
		}, nil
	}

	orderResp, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		PaymentMethod: protoPaymentMethod,
	})
	if err != nil {
		log.ErrorContext(ctx, "failed to process payment", "error", err)
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	// 4. Обновить статус на PAID и сохранить transaction_uuid

	transactionUUID, err := uuid.Parse(orderResp.GetTransactionUuid())
	if err != nil {
		log.ErrorContext(ctx, "invalid transaction id from payment service", "error", err)
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	paymentMethod := PaymentMethod(req.PaymentMethod)

	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	order.Status = OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethod

	// Обновление информации в Order

	h.store.orders[params.OrderUUID] = order

	log.InfoContext(ctx, "заказ оплачен", "order_uuid", params.OrderUUID, "transaction_uuid", transactionUUID)

	// 5. Вернуть transaction_uuid
	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel
func (h *handler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	log := slog.With("op", "order/pkg/handler/CancelOrder")

	// 1. Найти заказ в store

	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// Если не найден — вернуть 404

	if !ok {
		log.ErrorContext(ctx, "заказ не найден", "order_uuid", params.OrderUUID)
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("заказ %s не найден", params.OrderUUID),
		}, nil
	}
	// 2. Проверить статус == PENDING_PAYMENT

	if order.Status != OrderStatusPendingPayment {
		log.ErrorContext(ctx, "заказ не в статусе PENDING_PAYMENT", "order_uuid", params.OrderUUID, "status", order.Status)
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("заказ %s не в статусе PENDING_PAYMENT, его статус %s", params.OrderUUID, order.Status),
		}, nil
	}

	// 3. Обновить статус на CANCELLED

	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	order.Status = OrderStatusCancelled
	h.store.orders[params.OrderUUID] = order

	log.InfoContext(ctx, "заказ отменён", "order_uuid", params.OrderUUID)

	// 4. Вернуть success
	return &orderv1.CancelOrderResponse{}, nil
}
