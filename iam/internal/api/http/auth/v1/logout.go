package v1

import (
	"context"

	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// Logout завершает сессию, переданную в заголовке X-Session-UUID
func (a *api) Logout(ctx context.Context, params iamv1.LogoutParams) (iamv1.LogoutRes, error) {
	if err := a.sessionService.Logout(ctx, params.XSessionUUID); err != nil {
		return nil, err
	}

	return &iamv1.LogoutNoContent{}, nil
}
