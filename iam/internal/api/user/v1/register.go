package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/iam/internal/api/converter"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	op := "iam/internal/api/user/v1/RegisterUser"
	log := slog.With("op", op)
	userInput, err := converter.RegisterRequestToInput(req)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при конвертации запроса в input", "err", err)
		return nil, err
	}
	registerUuid, err := a.userService.Register(ctx, userInput)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при регистрации пользователя", "err", err)
		return nil, err
	}
	log.InfoContext(ctx, "пользователь успешно зарегистрирован", "user_uuid", registerUuid.String())
	return &userv1.RegisterResponse{
		UserUuid: registerUuid.String(),
	}, nil
}
