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

func (s *service) Register(ctx context.Context, input input.RegisterInput) (uuid.UUID, error) {
	op := "iam/internal/service/iam/Register"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		log.ErrorContext(ctx, "ошибка хэширования пароля", "error", err)
		return uuid.Nil, err
	}
	userModel := model.User{
		UUID:         uuid.New(),
		Login:        input.Login,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	err = s.userRepo.Create(ctx, userModel)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			log.ErrorContext(ctx, "пользователь уже существует", "login", input.Login)
			return uuid.Nil, err
		}
		log.ErrorContext(ctx, "ошибка создания пользователя", "error", err)
		return uuid.Nil, fmt.Errorf("создать пользователя: %w", err)
	}
	log.InfoContext(ctx, "пользователь успешно зарегистрирован", "uuid", userModel.UUID, "login", input.Login)
	return userModel.UUID, nil
}
