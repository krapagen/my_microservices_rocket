package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/iam/internal/api/converter"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

func (a *api) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	op := "iam/api/auth/v1/Logout"
	log := slog.With("op", op)
	sessionUUID, err := converter.LogoutRequestToInput(req)
	if err != nil {
		log.ErrorContext(ctx, "не удалось конвертировать запрос на выход", "error", err)
		return nil, err
	}
	err = a.sessionService.Logout(ctx, sessionUUID)
	if err != nil {
		log.ErrorContext(ctx, "не удалось выполнить выход", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "успешный выход")
	return &authv1.LogoutResponse{}, nil
}
