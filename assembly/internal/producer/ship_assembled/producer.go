package ship_assembled

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
	eventsv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	assembledProducer Producer
}

func NewService(assembledProducer Producer) ProducerService {
	return &service{
		assembledProducer: assembledProducer,
	}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {

	msg := &eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID.String(),
		OrderUuid:    event.OrderUUID.String(),
		UserUuid:     event.UserUUID.String(),
		BuildTimeSec: int64(event.BuildTime.Seconds()),
		AssembledAt:  timestamppb.New(event.AssembledAt),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось сериализовать ShipAssembled", "error", err)
		return err
	}

	return s.assembledProducer.Send(ctx, &kafka.Message{
		Key:   []byte(event.EventUUID.String()),
		Value: payload,
	})
}
