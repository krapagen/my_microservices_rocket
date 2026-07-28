package assembly_consumer

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	eventsv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/events/v1"
)

func decodeShipAssembled(data []byte) (model.ShipAssembled, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembled{}, fmt.Errorf("не удалось десериализовать protobuf: %w", err)
	}

	eventUUID, err := uuid.Parse(pb.GetEventUuid())
	if err != nil {
		return model.ShipAssembled{}, fmt.Errorf("invalid event_uuid: %w", err)
	}
	orderUUID, err := uuid.Parse(pb.GetOrderUuid())
	if err != nil {
		return model.ShipAssembled{}, fmt.Errorf("invalid order_uuid: %w", err)
	}
	userUUID, err := uuid.Parse(pb.GetUserUuid())
	if err != nil {
		return model.ShipAssembled{}, fmt.Errorf("invalid user_uuid: %w", err)
	}

	return model.ShipAssembled{
		EventUUID: eventUUID,
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
	}, nil
}
