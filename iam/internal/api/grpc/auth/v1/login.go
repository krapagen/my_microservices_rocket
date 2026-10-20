package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/iam/internal/api/converter"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	op := "iam/api/auth/v1/Login"
	log := slog.With("op", op)
	input, err := converter.LoginRequestToInput(req)
	if err != nil {
		log.ErrorContext(ctx, "не удалось конвертировать запрос на вход", "error", err)
		return nil, fmt.Errorf("конвертировать запрос на вход: %w", err)
	}
	uuid, err := a.sessionService.Login(ctx, input)
	if err != nil {
		log.ErrorContext(ctx, "не удалось залогинить пользователя", "login", input.Login, "error", err)
		return nil, fmt.Errorf("залогинить пользователя: %w", err)
	}
	log.InfoContext(ctx, "пользователь успешно залогинен", "login", input.Login)
	return &authv1.LoginResponse{
		SessionUuid: uuid.String(),
	}, nil
}
