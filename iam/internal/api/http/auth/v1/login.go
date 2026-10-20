package v1

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// Login аутентифицирует пользователя и создаёт сессию
func (a *api) Login(ctx context.Context, req *iamv1.LoginRequest) (iamv1.LoginRes, error) {
	sessionUUID, err := a.sessionService.Login(ctx, input.LoginInput{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &iamv1.LoginResponse{
		SessionUUID: sessionUUID,
	}, nil
}
