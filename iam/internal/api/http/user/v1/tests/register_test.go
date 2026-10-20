package tests

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

func (s *UserAPISuite) TestRegister_Success() {
	userUUID := uuid.New()
	req := &iamv1.RegisterRequest{
		Login:    "newuser",
		Password: "password123",
	}

	s.userService.EXPECT().Register(s.ctx, input.RegisterInput{
		Login:    "newuser",
		Password: "password123",
	}).Return(userUUID, nil)

	res, err := s.api.Register(s.ctx, req)

	s.NoError(err)
	s.Require().NotNil(res)
	resp, ok := res.(*iamv1.RegisterResponse)
	s.Require().True(ok)
	s.Equal(userUUID, resp.GetUserUUID())
}

func (s *UserAPISuite) TestRegister_EmptyLogin() {
	req := &iamv1.RegisterRequest{
		Login:    "",
		Password: "password123",
	}
	expectedInput := input.RegisterInput{Login: "", Password: "password123"}

	s.userService.EXPECT().Register(s.ctx, expectedInput).Return(uuid.Nil, errs.ErrInvalidLogin)

	res, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidLogin)
	s.Nil(res)
}

func (s *UserAPISuite) TestRegister_WeakPassword() {
	req := &iamv1.RegisterRequest{
		Login:    "newuser",
		Password: "short",
	}
	expectedInput := input.RegisterInput{Login: "newuser", Password: "short"}

	s.userService.EXPECT().Register(s.ctx, expectedInput).Return(uuid.Nil, errs.ErrWeakPassword)

	res, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrWeakPassword)
	s.Nil(res)
}

func (s *UserAPISuite) TestRegister_UserAlreadyExists() {
	req := &iamv1.RegisterRequest{
		Login:    "newuser",
		Password: "password123",
	}
	expectedInput := input.RegisterInput{Login: "newuser", Password: "password123"}

	s.userService.EXPECT().Register(s.ctx, expectedInput).Return(uuid.Nil, errs.ErrUserAlreadyExists)

	res, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrUserAlreadyExists)
	s.Nil(res)
}

func (s *UserAPISuite) TestRegister_ServiceError() {
	svcErr := errors.New("service error")
	req := &iamv1.RegisterRequest{
		Login:    "newuser",
		Password: "password123",
	}
	expectedInput := input.RegisterInput{Login: "newuser", Password: "password123"}

	s.userService.EXPECT().Register(s.ctx, expectedInput).Return(uuid.Nil, svcErr)

	res, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(res)
}
