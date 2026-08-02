package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/iam/internal/api/converter"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	op := "iam/internal/api/user/v1/GetUser"
	log := slog.With("op", op)
	userUUID, err := converter.GetUserRequestToInput(req)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при конвертации запроса в input", "err", err)
		return nil, err
	}
	user, err := a.userService.GetUser(ctx, userUUID)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при получении пользователя", "err", err)
		return nil, err
	}
	userDto, err := converter.UserToDto(user)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при конвертации пользователя в dto", "err", err)
		return nil, err
	}
	log.InfoContext(ctx, "пользователь успешно получен", "user_uuid", userUUID.String())
	return &userv1.GetUserResponse{
		User: userDto,
	}, nil
}
