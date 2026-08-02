package v1

import userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"

type api struct {
	userv1.UnimplementedUserServiceServer
	userService UserService
}

func New(userService UserService) *api {
	return &api{
		userService: userService,
	}
}
