package tests

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

func (s *AuthAPISuite) TestLogin_Success() {
	login := "validuser"
	password := "strongpass"
	sessionUUID := uuid.New()
	req := &authv1.LoginRequest{Login: login, Password: password}
	expectedInput := input.LoginInput{Login: login, Password: password}

	s.sessionService.EXPECT().Login(s.ctx, expectedInput).Return(sessionUUID, nil)

	resp, err := s.api.Login(s.ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(sessionUUID.String(), resp.GetSessionUuid())
}

func (s *AuthAPISuite) TestLogin_InvalidLogin() {
	req := &authv1.LoginRequest{Login: "", Password: "strongpass"}

	resp, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidLogin)
	s.Nil(resp)
	s.sessionService.AssertNotCalled(s.T(), "Login")
}

func (s *AuthAPISuite) TestLogin_EmptyPassword() {
	req := &authv1.LoginRequest{Login: "validuser", Password: ""}

	resp, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptyCredential)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestLogin_WeakPassword() {
	req := &authv1.LoginRequest{Login: "validuser", Password: "short"}

	resp, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrWeakPassword)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestLogin_ServiceError() {
	login := "validuser"
	password := "strongpass"
	req := &authv1.LoginRequest{Login: login, Password: password}
	expectedInput := input.LoginInput{Login: login, Password: password}
	svcErr := errors.New("service error")

	s.sessionService.EXPECT().Login(s.ctx, expectedInput).Return(uuid.Nil, svcErr)

	resp, err := s.api.Login(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}
