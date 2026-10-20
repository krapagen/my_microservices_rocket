package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

// SessionService — сервисный слой сессий (Redis)
type SessionService interface {
	Login(ctx context.Context, input input.LoginInput) (uuid.UUID, error)
	Logout(ctx context.Context, session uuid.UUID) error
}
