package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	userv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/user/v1"
	"github.com/krapagen/my_microservices_rocket/iam/internal/api/user/v1/mocks"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

type UserAPISuite struct {
	suite.Suite
	ctx         context.Context
	userService *mocks.UserService
	api         userv1.UserServiceServer
}

func (s *UserAPISuite) SetupTest() {
	s.ctx = context.Background()
	s.userService = mocks.NewUserService(s.T())
	s.api = userv1api.New(s.userService)
}

func TestUserAPISuite(t *testing.T) {
	suite.Run(t, new(UserAPISuite))
}
