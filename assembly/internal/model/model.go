package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderPaid struct {
	EventUUID uuid.UUID
	OrderUUID uuid.UUID
	UserUUID  uuid.UUID
}

type ShipAssembled struct {
	EventUUID   uuid.UUID
	OrderUUID   uuid.UUID
	UserUUID    uuid.UUID
	BuildTime   time.Duration
	AssembledAt time.Time
}
