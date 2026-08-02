package iam

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	op := "iam/internal/service/iam/GetUser"
	log := slog.With("op", op)
	if userUUID == uuid.Nil {
		log.ErrorContext(ctx, "user_uuid обязателен")
		return model.User{}, errs.ErrInvalidUUID
	}
	user, err := s.userRepo.GetByUUID(ctx, userUUID.String())
	if err != nil {
		log.ErrorContext(ctx, "ошибка получения пользователя", "uuid", userUUID, "error", err)
		return model.User{}, err
	}
	log.InfoContext(ctx, "пользователь получен", "uuid", userUUID)
	return user, nil
}
