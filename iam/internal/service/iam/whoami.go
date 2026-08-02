package iam

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
)

// WhoAmI возвращает сессию и пользователя по session_uuid.
func (s *service) WhoAmI(ctx context.Context, sessionUuid uuid.UUID) (model.Session, model.User, error) {
	op := "iam/internal/service/iam/WhoAmI"
	log := slog.With("op", op)
	if sessionUuid == uuid.Nil {
		log.ErrorContext(ctx, "session_uuid обязателен")
		return model.Session{}, model.User{}, errs.ErrEmptySessionID
	}
	session, err := s.sessionRepo.Get(ctx, sessionUuid.String())
	if err != nil {
		if errors.Is(err, errs.ErrSessionNotFound) {
			log.ErrorContext(ctx, "сессия не найдена или истекла", "session_uuid", sessionUuid, "error", err)
			return model.Session{}, model.User{}, errs.ErrSessionNotFound
		}
		log.ErrorContext(ctx, "не удалось получить сессию", "session_uuid", sessionUuid, "error", err)
		return model.Session{}, model.User{}, fmt.Errorf("получить сессию: %w", err)
	}
	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID.String())
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.ErrorContext(ctx, "пользователь не найден", "user_uuid", session.UserUUID, "error", err)
			return model.Session{}, model.User{}, errs.ErrUserNotFound
		}
		log.ErrorContext(ctx, "не удалось получить пользователя", "user_uuid", session.UserUUID, "error", err)
		return model.Session{}, model.User{}, fmt.Errorf("получить пользователя: %w", err)
	}
	return session, user, nil
}
