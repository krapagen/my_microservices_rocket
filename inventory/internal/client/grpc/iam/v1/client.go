package v1

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
	commonv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/common/v1"
)

// Client — клиент IAM-сервиса для inventory.
type Client interface {
	// Whoami проверяет сессию и возвращает данные пользователя и сессии.
	Whoami(ctx context.Context, sessionUUID string) (*commonv1.User, *commonv1.Session, error)
}

type client struct {
	authClient authv1.AuthServiceClient
}

// New создаёт обёртку над gRPC-клиентом IAM AuthService.
func New(conn grpc.ClientConnInterface) Client {
	return &client{authClient: authv1.NewAuthServiceClient(conn)}
}

// NewFromClient создаёт обёртку из уже инициализированного auth-клиента.
func NewFromClient(authClient authv1.AuthServiceClient) Client {
	return &client{authClient: authClient}
}

func (c *client) Whoami(ctx context.Context, sessionUUID string) (*commonv1.User, *commonv1.Session, error) {
	op := "inventory/internal/client/grpc/iam/v1/Whoami"
	log := slog.With("op", op)

	log.InfoContext(ctx, "Проверка сессии в IAM", "sessionUUID", sessionUUID)

	resp, err := c.authClient.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	if err != nil {
		log.ErrorContext(ctx, "Ошибка вызова IAM.Whoami", "error", err)
		switch status.Code(err) {
		case codes.NotFound:
			return nil, nil, fmt.Errorf("сессия не найдена: %w", err)
		case codes.Unauthenticated:
			return nil, nil, fmt.Errorf("недействительная сессия: %w", err)
		default:
			return nil, nil, fmt.Errorf("ошибка IAM: %w", err)
		}
	}

	log.InfoContext(ctx, "Сессия успешно проверена", "userUUID", resp.GetUser().GetUuid())
	return resp.GetUser(), resp.GetSession(), nil
}
