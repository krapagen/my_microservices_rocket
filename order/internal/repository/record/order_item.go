package record

import "github.com/google/uuid"

type OrderItem struct {
	OrderUUID uuid.UUID `db:"order_uuid"`
	PartUUID  uuid.UUID `db:"part_uuid"`
	PartType  string    `db:"part_type"`
	Price     int64     `db:"price"`
}
