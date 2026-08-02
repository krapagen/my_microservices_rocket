package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/redis_view"
)

func ToModelSession(session redis_view.SessionView) (model.Session, error) {
	if session.UUID == "" {
		return model.Session{}, errs.ErrSessionNotFound
	}
	uuidSession, err := uuid.Parse(session.UUID)
	if err != nil {
		return model.Session{}, fmt.Errorf("ошибка парсинга UUID сессии: %w", errs.ErrInvalidUUID)
	}
	userUUIDSession, err := uuid.Parse(session.UserUUID)
	if err != nil {
		return model.Session{}, fmt.Errorf("ошибка парсинга UUID пользователя: %w", errs.ErrInvalidUUID)
	}
	modelSession := model.Session{
		UUID:      uuidSession,
		UserUUID:  userUUIDSession,
		Login:     session.Login,
		CreatedAt: time.Unix(0, session.CreatedAtNs),
		ExpiresAt: time.Unix(0, session.ExpiresAtNs),
	}
	if session.UpdatedAtNs != nil {
		modelSession.UpdatedAt = new(time.Unix(0, *session.UpdatedAtNs))
	}
	return modelSession, nil
}

func ToSessionView(session model.Session) (redis_view.SessionView, error) {
	if session.UUID == uuid.Nil || session.UserUUID == uuid.Nil {
		return redis_view.SessionView{}, errs.ErrInvalidUUID
	}
	redisSession := redis_view.SessionView{
		UUID:        session.UUID.String(),
		UserUUID:    session.UserUUID.String(),
		Login:       session.Login,
		CreatedAtNs: session.CreatedAt.UnixNano(),
		ExpiresAtNs: session.ExpiresAt.UnixNano(),
	}
	if session.UpdatedAt != nil {
		redisSession.UpdatedAtNs = new(session.UpdatedAt.UnixNano())
	}
	return redisSession, nil
}
