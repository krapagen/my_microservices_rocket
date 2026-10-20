package v1

import (
	"context"

	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// AuthHandler — HTTP-обработчики аутентификации
type AuthHandler interface {
	Login(ctx context.Context, req *iamv1.LoginRequest) (iamv1.LoginRes, error)
	Logout(ctx context.Context, params iamv1.LogoutParams) (iamv1.LogoutRes, error)
}

// api реализует часть сгенерированного ogen-интерфейса iamv1.Handler
type api struct {
	sessionService SessionService
}

// New создаёт HTTP API-обработчик аутентификации
func New(sessionService SessionService) AuthHandler {
	return &api{
		sessionService: sessionService,
	}
}
