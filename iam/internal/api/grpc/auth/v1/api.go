package v1

import authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"

type api struct {
	authv1.UnimplementedAuthServiceServer
	sessionService SessionService
}

func New(sessionService SessionService) *api {
	return &api{
		sessionService: sessionService,
	}
}
