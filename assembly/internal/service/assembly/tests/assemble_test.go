package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/model"
	assemblysvc "github.com/krapagen/my_microservices_rocket/assembly/internal/service/assembly"
)

func TestAssemble_Success(t *testing.T) {
	buildTime := 0 * time.Second
	svc := assemblysvc.NewService(buildTime)

	in := model.OrderPaid{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
	}

	got, err := svc.Assemble(context.Background(), in)
	if err != nil {
		t.Fatalf("Assemble() error = %v", err)
	}

	if got.EventUUID == uuid.Nil {
		t.Fatal("Assemble() returned empty EventUUID")
	}
	if got.OrderUUID != in.OrderUUID {
		t.Fatalf("Assemble() OrderUUID = %s, want %s", got.OrderUUID, in.OrderUUID)
	}
	if got.UserUUID != in.UserUUID {
		t.Fatalf("Assemble() UserUUID = %s, want %s", got.UserUUID, in.UserUUID)
	}
	if got.BuildTime != buildTime {
		t.Fatalf("Assemble() BuildTime = %v, want %v", got.BuildTime, buildTime)
	}
	if got.AssembledAt.IsZero() {
		t.Fatal("Assemble() AssembledAt is zero")
	}
}

func TestAssemble_ContextCanceled(t *testing.T) {
	svc := assemblysvc.NewService(2 * time.Second)

	in := model.OrderPaid{
		EventUUID: uuid.New(),
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Assemble(ctx, in)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Assemble() error = %v, want %v", err, context.Canceled)
	}
}
