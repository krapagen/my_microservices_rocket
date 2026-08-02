package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/record"
)

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	return r.getBy(
		ctx,
		"iam/internal/repository/user/GetByLogin",
		"SELECT uuid, login, password_hash, created_at, updated_at FROM users WHERE login = $1",
		"login",
		login,
	)
}

func (r *repository) getBy(ctx context.Context, op, query, argKey, arg string) (model.User, error) {
	log := slog.With("op", op)
	var user record.UserRecord
	err := r.getter.DefaultTrOrDB(ctx, r.pool).
		QueryRow(ctx, query, arg).
		Scan(&user.UUID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.ErrorContext(ctx, "пользователь не найден", argKey, arg)
			return model.User{}, errs.ErrUserNotFound
		}
		log.ErrorContext(ctx, "ошибка получения пользователя", argKey, arg, "error", err)
		return model.User{}, err
	}
	log.InfoContext(ctx, "пользователь найден", argKey, arg)
	modelUser, err := converter.ToModelUser(user)
	if err != nil {
		log.ErrorContext(ctx, "ошибка конвертации пользователя", argKey, arg, "error", err)
		return model.User{}, err
	}
	return modelUser, nil
}
