package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	shipassembled "github.com/krapagen/my_microservices_rocket/assembly/internal/producer/ship_assembled"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
	eventsv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/events/v1"
)

type fakeKafkaProducer struct {
	sentMsg *kafka.Message
	err     error
}

func (f *fakeKafkaProducer) Send(_ context.Context, msg *kafka.Message) error {
	f.sentMsg = msg
	return f.err
}

func TestProduceShipAssembled_Success(t *testing.T) {
	fake := &fakeKafkaProducer{}
	svc := shipassembled.NewService(fake)

	event := model.ShipAssembled{
		EventUUID:   uuid.New(),
		OrderUUID:   uuid.New(),
		UserUUID:    uuid.New(),
		BuildTime:   7 * time.Second,
		AssembledAt: time.Now().UTC(),
	}

	err := svc.ProduceShipAssembled(context.Background(), event)
	if err != nil {
		t.Fatalf("ProduceShipAssembled() error = %v", err)
	}
	if fake.sentMsg == nil {
		t.Fatal("producer.Send() was not called")
	}
	if string(fake.sentMsg.Key) != event.EventUUID.String() {
		t.Fatalf("message key = %q, want %q", string(fake.sentMsg.Key), event.EventUUID.String())
	}

	var pb eventsv1.ShipAssembled
	if err = proto.Unmarshal(fake.sentMsg.Value, &pb); err != nil {
		t.Fatalf("proto.Unmarshal() error = %v", err)
	}

	if pb.GetEventUuid() != event.EventUUID.String() {
		t.Fatalf("event_uuid = %q, want %q", pb.GetEventUuid(), event.EventUUID.String())
	}
	if pb.GetOrderUuid() != event.OrderUUID.String() {
		t.Fatalf("order_uuid = %q, want %q", pb.GetOrderUuid(), event.OrderUUID.String())
	}
	if pb.GetUserUuid() != event.UserUUID.String() {
		t.Fatalf("user_uuid = %q, want %q", pb.GetUserUuid(), event.UserUUID.String())
	}
	if pb.GetBuildTimeSec() != 7 {
		t.Fatalf("build_time_sec = %d, want %d", pb.GetBuildTimeSec(), 7)
	}
	if pb.GetAssembledAt() == nil {
		t.Fatal("assembled_at is nil")
	}
}

func TestProduceShipAssembled_SendError(t *testing.T) {
	sendErr := errors.New("send failed")
	fake := &fakeKafkaProducer{err: sendErr}
	svc := shipassembled.NewService(fake)

	event := model.ShipAssembled{
		EventUUID:   uuid.New(),
		OrderUUID:   uuid.New(),
		UserUUID:    uuid.New(),
		BuildTime:   1 * time.Second,
		AssembledAt: time.Now().UTC(),
	}

	err := svc.ProduceShipAssembled(context.Background(), event)
	if !errors.Is(err, sendErr) {
		t.Fatalf("ProduceShipAssembled() error = %v, want %v", err, sendErr)
	}
}
