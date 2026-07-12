package record

import (
	"time"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

type PartType int32

const (
	PartType_PART_TYPE_UNSPECIFIED PartType = 0
	PartType_PART_TYPE_HULL        PartType = 1
	PartType_PART_TYPE_ENGINE      PartType = 2
	PartType_PART_TYPE_SHIELD      PartType = 3
	PartType_PART_TYPE_WEAPON      PartType = 4
)
