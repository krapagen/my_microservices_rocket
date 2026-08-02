package iam

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
)

func (s *service) Logout(ctx context.Context, session uuid.UUID) error {
	op := "iam/internal/service/iam/Logout"
	log := slog.With("op", op)
	if session == uuid.Nil {
		log.ErrorContext(ctx, "session_uuid обязателен")
		return errs.ErrEmptySessionID
	}
	err := s.sessionRepo.Delete(ctx, session.String())
	if err != nil {
		log.ErrorContext(ctx, "не удалось удалить сессию", "error", err)
		return err
	}
	log.InfoContext(ctx, "сессия успешно удалена", "session_uuid", session)
	return nil
}
