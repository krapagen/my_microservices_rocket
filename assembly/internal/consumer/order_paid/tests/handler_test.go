package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	orderpaid "github.com/krapagen/my_microservices_rocket/assembly/internal/consumer/order_paid"
	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/kafka"
	eventsv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/events/v1"
)

type fakeConsumer struct {
	msg        kafka.Message
	consumeErr error
}

func (f *fakeConsumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	if f.consumeErr != nil {
		return f.consumeErr
	}

	return handler(ctx, f.msg)
}

type fakeAssembler struct {
	result model.ShipAssembled
	err    error
	called bool
	input  model.OrderPaid
}

func (f *fakeAssembler) Assemble(_ context.Context, in model.OrderPaid) (model.ShipAssembled, error) {
	f.called = true
	f.input = in
	return f.result, f.err
}

type fakeShipAssembledProducer struct {
	err    error
	called bool
	input  model.ShipAssembled
}

func (f *fakeShipAssembledProducer) ProduceShipAssembled(_ context.Context, event model.ShipAssembled) error {
	f.called = true
	f.input = event
	return f.err
}

func TestOrderPaidHandler_Success(t *testing.T) {
	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	payload, err := proto.Marshal(&eventsv1.OrderPaid{
		EventUuid: eventUUID.String(),
		OrderUuid: orderUUID.String(),
		UserUuid:  userUUID.String(),
	})
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	expectedAssembled := model.ShipAssembled{
		EventUUID:   uuid.New(),
		OrderUUID:   orderUUID,
		UserUUID:    userUUID,
		BuildTime:   5 * time.Second,
		AssembledAt: time.Now().UTC(),
	}

	consumer := &fakeConsumer{
		msg: kafka.Message{Value: payload},
	}
	assembler := &fakeAssembler{result: expectedAssembled}
	producer := &fakeShipAssembledProducer{}

	svc := orderpaid.NewService(consumer, assembler, producer)
	err = svc.RunConsumer(context.Background())
	if err != nil {
		t.Fatalf("RunConsumer() error = %v", err)
	}

	if !assembler.called {
		t.Fatal("assembler.Assemble() was not called")
	}
	if assembler.input.EventUUID != eventUUID || assembler.input.OrderUUID != orderUUID || assembler.input.UserUUID != userUUID {
		t.Fatal("assembler.Assemble() got incorrect input")
	}
	if !producer.called {
		t.Fatal("producer.ProduceShipAssembled() was not called")
	}
	if producer.input != expectedAssembled {
		t.Fatalf("producer input mismatch: got %+v, want %+v", producer.input, expectedAssembled)
	}
}

func TestOrderPaidHandler_DecodeError(t *testing.T) {
	consumer := &fakeConsumer{
		msg: kafka.Message{Value: []byte("not-protobuf")},
	}
	assembler := &fakeAssembler{}
	producer := &fakeShipAssembledProducer{}

	svc := orderpaid.NewService(consumer, assembler, producer)
	err := svc.RunConsumer(context.Background())
	if err == nil {
		t.Fatal("RunConsumer() error = nil, want non-nil")
	}
	if assembler.called {
		t.Fatal("assembler.Assemble() should not be called on decode error")
	}
	if producer.called {
		t.Fatal("producer.ProduceShipAssembled() should not be called on decode error")
	}
}

func TestOrderPaidHandler_AssembleError(t *testing.T) {
	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	payload, err := proto.Marshal(&eventsv1.OrderPaid{
		EventUuid: eventUUID.String(),
		OrderUuid: orderUUID.String(),
		UserUuid:  userUUID.String(),
	})
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	assembleErr := errors.New("assemble failed")
	consumer := &fakeConsumer{msg: kafka.Message{Value: payload}}
	assembler := &fakeAssembler{err: assembleErr}
	producer := &fakeShipAssembledProducer{}

	svc := orderpaid.NewService(consumer, assembler, producer)
	err = svc.RunConsumer(context.Background())
	if !errors.Is(err, assembleErr) {
		t.Fatalf("RunConsumer() error = %v, want %v", err, assembleErr)
	}
	if producer.called {
		t.Fatal("producer.ProduceShipAssembled() should not be called on assemble error")
	}
}

func TestOrderPaidHandler_ProduceError(t *testing.T) {
	eventUUID := uuid.New()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	payload, err := proto.Marshal(&eventsv1.OrderPaid{
		EventUuid: eventUUID.String(),
		OrderUuid: orderUUID.String(),
		UserUuid:  userUUID.String(),
	})
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	produceErr := errors.New("produce failed")
	consumer := &fakeConsumer{msg: kafka.Message{Value: payload}}
	assembler := &fakeAssembler{result: model.ShipAssembled{
		EventUUID:   uuid.New(),
		OrderUUID:   orderUUID,
		UserUUID:    userUUID,
		BuildTime:   1 * time.Second,
		AssembledAt: time.Now().UTC(),
	}}
	producer := &fakeShipAssembledProducer{err: produceErr}

	svc := orderpaid.NewService(consumer, assembler, producer)
	err = svc.RunConsumer(context.Background())
	if !errors.Is(err, produceErr) {
		t.Fatalf("RunConsumer() error = %v, want %v", err, produceErr)
	}
}
