package tests

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	commonv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/common/v1"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

func (s *UserAPISuite) TestRegister_Success() {
	userUUID := uuid.New()
	req := &userv1.RegisterRequest{
		Info: &userv1.UserRegistrationInfo{
			Info: &commonv1.UserInfo{
				Login: "newuser",
			},
			Password: "password123",
		},
	}

	s.userService.EXPECT().Register(s.ctx, input.RegisterInput{
		Login:    "newuser",
		Password: "password123",
	}).Return(userUUID, nil)

	resp, err := s.api.Register(s.ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(userUUID.String(), resp.GetUserUuid())
}

func (s *UserAPISuite) TestRegister_EmptyLogin() {
	req := &userv1.RegisterRequest{
		Info: &userv1.UserRegistrationInfo{
			Info: &commonv1.UserInfo{
				Login: "",
			},
			Password: "password123",
		},
	}

	resp, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidLogin)
	s.Nil(resp)
}

func (s *UserAPISuite) TestRegister_EmptyPassword() {
	req := &userv1.RegisterRequest{
		Info: &userv1.UserRegistrationInfo{
			Info: &commonv1.UserInfo{
				Login: "newuser",
			},
			Password: "",
		},
	}

	resp, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptyCredential)
	s.Nil(resp)
}

func (s *UserAPISuite) TestRegister_WeakPassword() {
	req := &userv1.RegisterRequest{
		Info: &userv1.UserRegistrationInfo{
			Info: &commonv1.UserInfo{
				Login: "newuser",
			},
			Password: "short",
		},
	}

	resp, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrWeakPassword)
	s.Nil(resp)
}

func (s *UserAPISuite) TestRegister_ServiceError() {
	svcErr := errors.New("service error")
	req := &userv1.RegisterRequest{
		Info: &userv1.UserRegistrationInfo{
			Info: &commonv1.UserInfo{
				Login: "newuser",
			},
			Password: "password123",
		},
	}

	s.userService.EXPECT().Register(s.ctx, input.RegisterInput{
		Login:    "newuser",
		Password: "password123",
	}).Return(uuid.Nil, svcErr)

	resp, err := s.api.Register(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}
