package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	userv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/user/v1"
	"github.com/krapagen/my_microservices_rocket/iam/internal/api/http/user/v1/mocks"
)

type UserAPISuite struct {
	suite.Suite
	ctx            context.Context
	userService    *mocks.UserService
	sessionService *mocks.SessionService
	api            userv1api.UserHandler
}

func (s *UserAPISuite) SetupTest() {
	s.ctx = context.Background()
	s.userService = mocks.NewUserService(s.T())
	s.sessionService = mocks.NewSessionService(s.T())
	s.api = userv1api.New(s.userService, s.sessionService)
}

func TestUserAPISuite(t *testing.T) {
	suite.Run(t, new(UserAPISuite))
}
