package tests

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

func (s *AuthAPISuite) TestLogout_Success() {
	sessionUUID := uuid.New()
	req := &authv1.LogoutRequest{SessionUuid: sessionUUID.String()}

	s.sessionService.EXPECT().Logout(s.ctx, sessionUUID).Return(nil)

	resp, err := s.api.Logout(s.ctx, req)

	s.NoError(err)
	s.NotNil(resp)
}

func (s *AuthAPISuite) TestLogout_EmptySession() {
	req := &authv1.LogoutRequest{SessionUuid: ""}

	resp, err := s.api.Logout(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptySessionID)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestLogout_InvalidUUID() {
	req := &authv1.LogoutRequest{SessionUuid: "not-a-uuid"}

	resp, err := s.api.Logout(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestLogout_ServiceError() {
	sessionUUID := uuid.New()
	req := &authv1.LogoutRequest{SessionUuid: sessionUUID.String()}
	svcErr := errors.New("service error")

	s.sessionService.EXPECT().Logout(s.ctx, sessionUUID).Return(svcErr)

	resp, err := s.api.Logout(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}
