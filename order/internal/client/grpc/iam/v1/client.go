package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

type Client struct {
	authClient authv1.AuthServiceClient
}

// New создаёт обёртку над auth-клиентом.
func New(authClient authv1.AuthServiceClient) *Client {
	return &Client{authClient: authClient}
}

func (c *Client) Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error) {
	op := "order/internal/client/grpc/iam/v1/Whoami"
	log := slog.With("op", op)

	resp, err := c.authClient.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	if err != nil {
		log.ErrorContext(ctx, "ошибка вызова IAM.Whoami", "error", err)
		return uuid.Nil, fmt.Errorf("проверить сессию: %w", err)
	}

	userUUID, err := uuid.Parse(resp.GetUser().GetUuid())
	if err != nil {
		log.ErrorContext(ctx, "некорректный UUID пользователя в ответе IAM", "error", err)
		return uuid.Nil, fmt.Errorf("распарсить UUID пользователя: %w", err)
	}

	return userUUID, nil
}
