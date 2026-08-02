package iam

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetByUUID(ctx context.Context, uuid string) (model.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, uuid string) (model.Session, error)
	Set(ctx context.Context, uuid string, session model.Session, ttl time.Duration) error
	Delete(ctx context.Context, uuid string) error
}

type Service interface {
	Register(ctx context.Context, input input.RegisterInput) (uuid.UUID, error)
	Login(ctx context.Context, input input.LoginInput) (uuid.UUID, error)
	Logout(ctx context.Context, session uuid.UUID) error
	WhoAmI(ctx context.Context, sessionUuid uuid.UUID) (model.Session, model.User, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}
