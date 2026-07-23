// internal/errors/errors.go
package errs

import "errors"

var (
	// Ошибки заказов
	ErrOrderItemNotFound    = errors.New("деталь заказа не найдена")
	ErrOrderNotFound        = errors.New("заказ не найден")
	ErrOrderAlreadyPaid     = errors.New("заказ уже оплачен")
	ErrOrderCancelled       = errors.New("заказ отменён")
	ErrOrderStatusIncorrect = errors.New("неверный статус заказа")

	// Ошибки деталей
	ErrPartNotFound         = errors.New("деталь не найдена")
	ErrOutOfStock           = errors.New("деталь отсутствует на складе")
	ErrMissingRequiredParts = errors.New("не указаны обязательные детали")

	// Ошибки валидации
	ErrInvalidUUID          = errors.New("неверный формат UUID")
	ErrInvalidPaymentMethod = errors.New("неверный метод оплаты")

	// Ошибки совместимости деталей
	ErrIncompatibleParts = errors.New("детали несовместимы")
	ErrPartTypeMismatch  = errors.New("тип детали не соответствует слоту")
)
