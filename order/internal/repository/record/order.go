package record

import (
	"time"
)

type Order struct {
	OrderUUID       string
	Items           []OrderItem
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}
type OrderItem struct {
	PartUUID string
	PartType PartType
	Price    int64
}

type PartType string

const (
	PartTypeHull   PartType = "HULL"
	PartTypeEngine PartType = "ENGINE"
	PartTypeShield PartType = "SHIELD"
	PartTypeWeapon PartType = "WEAPON"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
