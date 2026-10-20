package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

type SessionService interface {
	WhoAmI(ctx context.Context, sessionUuid uuid.UUID) (model.Session, model.User, error)
	Login(ctx context.Context, input input.LoginInput) (uuid.UUID, error)
	Logout(ctx context.Context, session uuid.UUID) error
}
