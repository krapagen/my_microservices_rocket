package record

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID            uuid.UUID  `db:"uuid"`
	Status          string     `db:"status"`
	TransactionUUID *uuid.UUID `db:"transaction_uuid"` // nullable → указатель
	PaymentMethod   *string    `db:"payment_method"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
