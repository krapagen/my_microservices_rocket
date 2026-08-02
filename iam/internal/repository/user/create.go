package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, user model.User) error {
	op := "iam/internal/repository/user/Create"
	log := slog.With("op", op)
	userRecord := converter.ToUserRecord(user)
	query := `INSERT INTO users (uuid, login, password_hash, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).
		Exec(ctx, query, userRecord.UUID, userRecord.Login, userRecord.PasswordHash, userRecord.CreatedAt, userRecord.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.ErrorContext(ctx, "пользователь с таким логином уже существует", "login", user.Login)
			return errs.ErrUserAlreadyExists
		}
		log.ErrorContext(ctx, "ошибка создания пользователя", "error", err)
		return err
	}
	log.InfoContext(ctx, "пользователь создан", "uuid", user.UUID)
	return nil
}
