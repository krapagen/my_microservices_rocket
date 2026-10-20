package v1

import (
	iamauthv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/auth/v1"
	iamuserv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/user/v1"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// api — композиция доменных HTTP-обработчиков (auth + user)
// в единый сгенерированный ogen-интерфейс iamv1.Handler
type api struct {
	iamauthv1api.AuthHandler
	iamuserv1api.UserHandler
}

// New создаёт HTTP API-обработчик IAM-сервиса
func New(authHandler iamauthv1api.AuthHandler, userHandler iamuserv1api.UserHandler) *api {
	return &api{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}
}

// compile-time проверка соответствия сгенерированному интерфейсу
var _ iamv1.Handler = (*api)(nil)
