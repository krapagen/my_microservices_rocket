package tests

import (
	"errors"

	"github.com/google/uuid"

	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// Проверяем, что композиция делегирует каждую операцию в нужный доменный handler

func (s *IamAPISuite) TestLogin_DelegatesToAuth() {
	sessionUUID := uuid.New()
	s.auth.loginRes = &iamv1.LoginResponse{SessionUUID: sessionUUID}
	req := &iamv1.LoginRequest{Login: "u", Password: "password123"}

	res, err := s.api.Login(s.ctx, req)

	s.NoError(err)
	s.True(s.auth.loginCalled)
	s.False(s.user.registerCalled)
	s.Equal(sessionUUID, res.(*iamv1.LoginResponse).GetSessionUUID())
}

func (s *IamAPISuite) TestLogin_PropagatesError() {
	svcErr := errors.New("auth error")
	s.auth.loginErr = svcErr
	req := &iamv1.LoginRequest{Login: "u", Password: "password123"}

	res, err := s.api.Login(s.ctx, req)

	s.ErrorIs(err, svcErr)
	s.Nil(res)
}

func (s *IamAPISuite) TestLogout_DelegatesToAuth() {
	s.auth.logoutRes = &iamv1.LogoutNoContent{}
	params := iamv1.LogoutParams{XSessionUUID: uuid.New()}

	res, err := s.api.Logout(s.ctx, params)

	s.NoError(err)
	s.True(s.auth.logoutCalled)
	s.NotNil(res)
}

func (s *IamAPISuite) TestRegister_DelegatesToUser() {
	userUUID := uuid.New()
	s.user.registerRes = &iamv1.RegisterResponse{UserUUID: userUUID}
	req := &iamv1.RegisterRequest{Login: "u", Password: "password123"}

	res, err := s.api.Register(s.ctx, req)

	s.NoError(err)
	s.True(s.user.registerCalled)
	s.False(s.auth.loginCalled)
	s.Equal(userUUID, res.(*iamv1.RegisterResponse).GetUserUUID())
}

func (s *IamAPISuite) TestGetUser_DelegatesToUser() {
	s.user.getUserRes = &iamv1.UserDto{UUID: uuid.New(), Login: "u"}
	params := iamv1.GetUserParams{UserUUID: uuid.New(), XSessionUUID: uuid.New()}

	res, err := s.api.GetUser(s.ctx, params)

	s.NoError(err)
	s.True(s.user.getUserCalled)
	s.NotNil(res)
}

func (s *IamAPISuite) TestGetUser_PropagatesError() {
	svcErr := errors.New("user error")
	s.user.getUserErr = svcErr
	params := iamv1.GetUserParams{UserUUID: uuid.New(), XSessionUUID: uuid.New()}

	res, err := s.api.GetUser(s.ctx, params)

	s.ErrorIs(err, svcErr)
	s.Nil(res)
}
