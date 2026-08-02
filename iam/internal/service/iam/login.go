package iam

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

func (s *service) Login(ctx context.Context, input input.LoginInput) (uuid.UUID, error) {
	op := "iam/internal/service/iam/Login"
	log := slog.With("op", op)
	if input.Login == "" {
		return uuid.Nil, errs.ErrInvalidLogin
	}
	if input.Password == "" {
		return uuid.Nil, errs.ErrEmptyCredential
	}
	if len(input.Password) < 8 {
		return uuid.Nil, errs.ErrWeakPassword
	}
	user, err := s.userRepo.GetByLogin(ctx, input.Login)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.ErrorContext(ctx, "пользователь не найден", "login", input.Login)
			return uuid.Nil, errs.ErrInvalidCredentials
		}
		log.ErrorContext(ctx, "не удалось получить пользователя", "login", input.Login, "error", err)
		return uuid.Nil, fmt.Errorf("получить пользователя: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		log.ErrorContext(ctx, "неверный пароль", "login", input.Login, "error", err)
		return uuid.Nil, errs.ErrInvalidCredentials
	}

	now := time.Now()
	session := model.Session{
		UUID:      uuid.New(),
		UserUUID:  user.UUID,
		Login:     user.Login,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
	}
	err = s.sessionRepo.Set(ctx, session.UUID.String(), session, s.sessionTTL)
	if err != nil {
		log.ErrorContext(ctx, "не удалось создать сессию", "login", input.Login, "error", err)
		return uuid.Nil, fmt.Errorf("создать сессию: %w", err)
	}
	log.InfoContext(ctx, "пользователь успешно вошел в систему", "uuid", user.UUID, "login", input.Login, "session_uuid", session.UUID)
	return session.UUID, nil
}
