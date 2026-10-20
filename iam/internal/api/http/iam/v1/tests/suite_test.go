package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	iamv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/iam/v1"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// authHandlerStub — тестовый стаб HTTP-обработчика аутентификации.
// Записывает факт вызова и возвращает заготовленный ответ.
type authHandlerStub struct {
	loginCalled  bool
	logoutCalled bool
	loginRes     iamv1.LoginRes
	loginErr     error
	logoutRes    iamv1.LogoutRes
	logoutErr    error
}

func (s *authHandlerStub) Login(_ context.Context, _ *iamv1.LoginRequest) (iamv1.LoginRes, error) {
	s.loginCalled = true
	return s.loginRes, s.loginErr
}

func (s *authHandlerStub) Logout(_ context.Context, _ iamv1.LogoutParams) (iamv1.LogoutRes, error) {
	s.logoutCalled = true
	return s.logoutRes, s.logoutErr
}

// userHandlerStub — тестовый стаб HTTP-обработчика пользователей.
type userHandlerStub struct {
	registerCalled bool
	getUserCalled  bool
	registerRes    iamv1.RegisterRes
	registerErr    error
	getUserRes     iamv1.GetUserRes
	getUserErr     error
}

func (s *userHandlerStub) Register(_ context.Context, _ *iamv1.RegisterRequest) (iamv1.RegisterRes, error) {
	s.registerCalled = true
	return s.registerRes, s.registerErr
}

func (s *userHandlerStub) GetUser(_ context.Context, _ iamv1.GetUserParams) (iamv1.GetUserRes, error) {
	s.getUserCalled = true
	return s.getUserRes, s.getUserErr
}

type IamAPISuite struct {
	suite.Suite
	ctx  context.Context
	auth *authHandlerStub
	user *userHandlerStub
	api  iamv1.Handler
}

func (s *IamAPISuite) SetupTest() {
	s.ctx = context.Background()
	s.auth = &authHandlerStub{}
	s.user = &userHandlerStub{}
	s.api = iamv1api.New(s.auth, s.user)
}

func TestIamAPISuite(t *testing.T) {
	suite.Run(t, new(IamAPISuite))
}
