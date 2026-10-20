package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

type UserService interface {
	Register(ctx context.Context, input input.RegisterInput) (uuid.UUID, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}
