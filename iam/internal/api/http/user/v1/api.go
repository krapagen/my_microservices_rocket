package v1

import (
	"context"

	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// UserHandler — HTTP-обработчики пользователей
type UserHandler interface {
	Register(ctx context.Context, req *iamv1.RegisterRequest) (iamv1.RegisterRes, error)
	GetUser(ctx context.Context, params iamv1.GetUserParams) (iamv1.GetUserRes, error)
}

// api реализует часть сгенерированного ogen-интерфейса iamv1.Handler
type api struct {
	userService    UserService
	sessionService SessionService
}

// New создаёт HTTP API-обработчик пользователей
func New(userService UserService, sessionService SessionService) UserHandler {
	return &api{
		userService:    userService,
		sessionService: sessionService,
	}
}
