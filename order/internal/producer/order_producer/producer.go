package order_producer

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
	eventsv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/events/v1"
)

type Producer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}

type ProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}

type service struct {
	orderPaidProducer Producer
}

func New(orderPaidProducer Producer) ProducerService {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error {
	msg := &eventsv1.OrderPaid{
		EventUuid: event.EventUUID.String(),
		OrderUuid: event.OrderUUID.String(),
		UserUuid:  event.UserUUID.String(),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		slog.ErrorContext(ctx, "не удалось сериализовать OrderPaid", "error", err)
		return err
	}

	return s.orderPaidProducer.Send(ctx, &kafka.Message{
		Key:   []byte(event.EventUUID.String()),
		Value: payload,
	})
}
