package v1

import (
	"context"
	"log/slog"

	"github.com/krapagen/my_microservices_rocket/iam/internal/api/converter"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	op := "iam/internal/api/auth/v1/WhoAmI"
	log := slog.With("op", op)
	sessionUUID, err := converter.WhoAmIRequestToInput(req)
	if err != nil {
		log.ErrorContext(ctx, "не удалось конвертировать запрос на получение информации о пользователе", "error", err)
		return nil, err
	}
	session, user, err := a.sessionService.WhoAmI(ctx, sessionUUID)
	if err != nil {
		log.ErrorContext(ctx, "не удалось получить информацию о пользователе", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "успешное получение информации о пользователе")
	sessionDto, err := converter.SessionToDto(session)
	if err != nil {
		log.ErrorContext(ctx, "не удалось конвертировать сессию в DTO", "error", err)
		return nil, err
	}
	userDto, err := converter.UserToDto(user)
	if err != nil {
		log.ErrorContext(ctx, "не удалось конвертировать пользователя в DTO", "error", err)
		return nil, err
	}
	log.InfoContext(ctx, "успешная конвертация сессии и пользователя в DTO")
	return &authv1.WhoamiResponse{
		Session: sessionDto,
		User:    userDto,
	}, nil
}
