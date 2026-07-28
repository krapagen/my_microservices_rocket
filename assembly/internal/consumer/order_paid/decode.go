package order_paid

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	eventsv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/events/v1"
)

func decodeOrderPaid(data []byte) (model.OrderPaid, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaid{}, fmt.Errorf("не удалось десериализовать protobuf: %w", err)
	}
	eventUUID, err := uuid.Parse(pb.EventUuid)
	if err != nil {
		return model.OrderPaid{}, fmt.Errorf("invalid event_uuid: %w", err)
	}
	orderUUID, err := uuid.Parse(pb.OrderUuid)
	if err != nil {
		return model.OrderPaid{}, fmt.Errorf("invalid order_uuid: %w", err)
	}
	userUUID, err := uuid.Parse(pb.UserUuid)
	if err != nil {
		return model.OrderPaid{}, fmt.Errorf("invalid user_uuid: %w", err)
	}
	return model.OrderPaid{
		EventUUID: eventUUID,
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
	}, nil
}
