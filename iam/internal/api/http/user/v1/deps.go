package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

// UserService — сервисный слой пользователей (PostgreSQL)
type UserService interface {
	Register(ctx context.Context, input input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}

// SessionService — сервисный слой сессий (Redis) для проверки авторизации
type SessionService interface {
	WhoAmI(ctx context.Context, sessionUuid uuid.UUID) (model.Session, model.User, error)
}
