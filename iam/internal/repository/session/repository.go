package session

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/redis_view"
)

const (
	cacheKeyPrefix = "session:"
)

type redisClient interface {
	HSet(ctx context.Context, key string, values ...any) *redis.IntCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type repository struct {
	client redisClient
}

func NewRepository(client redisClient) *repository {
	return &repository{
		client: client,
	}
}

func (r *repository) getCacheKey(uuid string) string {
	return cacheKeyPrefix + uuid
}

func (r *repository) Get(ctx context.Context, uuid string) (model.Session, error) {
	op := "iam/repository/session.Get"
	log := slog.With("op", op)
	cacheKey := r.getCacheKey(uuid)

	var sessionRedisView redis_view.SessionView
	err := r.client.HGetAll(ctx, cacheKey).Scan(&sessionRedisView)
	if err != nil {
		// redis.Nil от HGetAll не приходит — эту ошибку возвращают только
		// строковые команды (GET, HGET и т.п.) для отсутствующего ключа.
		// Здесь ветка нужна только на случай, если кто-то прокинет Nil
		// сверху (например, обёрткой клиента). Реальный "not found"
		// у HGetAll ловится проверкой на пустой UUID ниже.
		if errors.Is(err, redis.Nil) {
			log.ErrorContext(ctx, "сессия не найдена", "uuid", uuid)
			return model.Session{}, errs.ErrSessionNotFound
		}
		log.ErrorContext(ctx, "ошибка при получении сессии", "uuid", uuid, "err", err)
		return model.Session{}, err
	}
	if sessionRedisView.UUID == "" {
		log.ErrorContext(ctx, "сессия не найдена", "uuid", uuid)
		return model.Session{}, errs.ErrSessionNotFound
	}
	session, err := converter.ToModelSession(sessionRedisView)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при конвертации сессии", "uuid", uuid, "err", err)
		return model.Session{}, err
	}
	log.InfoContext(ctx, "сессия успешно получена", "uuid", uuid)
	return session, nil
}

func (r *repository) Set(ctx context.Context, uuid string, session model.Session, ttl time.Duration) error {
	op := "iam/repository/session.Set"
	log := slog.With("op", op)
	cacheKey := r.getCacheKey(uuid)
	sessionView, err := converter.ToSessionView(session)
	if err != nil {
		log.ErrorContext(ctx, "ошибка при конвертации сессии", "uuid", uuid, "err", err)
		return err
	}
	err = r.client.HSet(ctx, cacheKey, sessionView).Err()
	if err != nil {
		log.ErrorContext(ctx, "ошибка при установке сессии", "uuid", uuid, "err", err)
		return err
	}
	err = r.client.Expire(ctx, cacheKey, ttl).Err()
	if err != nil {
		log.ErrorContext(ctx, "ошибка при установке времени жизни сессии", "uuid", uuid, "err", err)
		return err
	}
	log.InfoContext(ctx, "сессия успешно установлена", "uuid", uuid)
	return nil
}

func (r *repository) Delete(ctx context.Context, uuid string) error {
	op := "iam/repository/session.Delete"
	log := slog.With("op", op)
	cacheKey := r.getCacheKey(uuid)
	err := r.client.Del(ctx, cacheKey).Err()
	if err != nil {
		log.ErrorContext(ctx, "ошибка при удалении сессии", "uuid", uuid, "err", err)
		return err
	}
	log.InfoContext(ctx, "сессия успешно удалена", "uuid", uuid)
	return nil
}
