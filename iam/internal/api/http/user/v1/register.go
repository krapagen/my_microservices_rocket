package v1

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// Register регистрирует нового пользователя
func (a *api) Register(ctx context.Context, req *iamv1.RegisterRequest) (iamv1.RegisterRes, error) {
	userUUID, err := a.userService.Register(ctx, input.RegisterInput{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &iamv1.RegisterResponse{
		UserUUID: userUUID,
	}, nil
}
