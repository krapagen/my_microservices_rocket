package v1

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/iam/internal/api/converter"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// GetUser возвращает пользователя по UUID.
// Требует активную сессию в заголовке X-Session-UUID — иначе 401
func (a *api) GetUser(ctx context.Context, params iamv1.GetUserParams) (iamv1.GetUserRes, error) {
	if _, _, err := a.sessionService.WhoAmI(ctx, params.XSessionUUID); err != nil {
		return nil, err
	}

	user, err := a.userService.GetUser(ctx, params.UserUUID)
	if err != nil {
		return nil, err
	}

	return converter.UserModelToHTTPDTO(user), nil
}
