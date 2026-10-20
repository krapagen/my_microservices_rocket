package tests

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

func (s *AuthAPISuite) TestLogin_Success() {
	login := "validuser"
	password := "strongpass"
	sessionUUID := uuid.New()
	req := &iamv1.LoginRequest{Login: login, Password: password}
	expectedInput := input.LoginInput{Login: login, Password: password}

	s.sessionService.EXPECT().Login(s.ctx, expectedInput).Return(sessionUUID, nil)

	res, err := s.api.Login(s.ctx, req)

	s.NoError(err)
	s.Require().NotNil(res)
	resp, ok := res.(*iamv1.LoginResponse)
	s.Require().True(ok)
	s.Equal(sessionUUID, resp.GetSessionUUID())
}

func (s *AuthAPISuite) TestLogin_EmptyCredentials() {
	req := &iamv1.LoginRequest{Login: "", Password: ""}
	expectedInput := input.LoginInput{}

	s.sessionService.EXPECT().Login(s.ctx, expectedInput).Return(uuid.Nil, errs.ErrEmptyCredential)

	res, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptyCredential)
	s.Nil(res)
}

func (s *AuthAPISuite) TestLogin_InvalidCredentials() {
	login := "validuser"
	password := "wrongpass"
	req := &iamv1.LoginRequest{Login: login, Password: password}
	expectedInput := input.LoginInput{Login: login, Password: password}

	s.sessionService.EXPECT().Login(s.ctx, expectedInput).Return(uuid.Nil, errs.ErrInvalidCredentials)

	res, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidCredentials)
	s.Nil(res)
}

func (s *AuthAPISuite) TestLogin_ServiceError() {
	login := "validuser"
	password := "strongpass"
	req := &iamv1.LoginRequest{Login: login, Password: password}
	expectedInput := input.LoginInput{Login: login, Password: password}
	svcErr := errors.New("service error")

	s.sessionService.EXPECT().Login(s.ctx, expectedInput).Return(uuid.Nil, svcErr)

	res, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(res)
}
