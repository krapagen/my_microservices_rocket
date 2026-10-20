package tests

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

func (s *AuthAPISuite) TestLogout_Success() {
	sessionUUID := uuid.New()
	params := iamv1.LogoutParams{XSessionUUID: sessionUUID}

	s.sessionService.EXPECT().Logout(s.ctx, sessionUUID).Return(nil)

	res, err := s.api.Logout(s.ctx, params)

	s.NoError(err)
	s.Require().NotNil(res)
	_, ok := res.(*iamv1.LogoutNoContent)
	s.True(ok)
}

func (s *AuthAPISuite) TestLogout_EmptySession() {
	params := iamv1.LogoutParams{XSessionUUID: uuid.Nil}

	s.sessionService.EXPECT().Logout(s.ctx, uuid.Nil).Return(errs.ErrEmptySessionID)

	res, err := s.api.Logout(s.ctx, params)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptySessionID)
	s.Nil(res)
}

func (s *AuthAPISuite) TestLogout_ServiceError() {
	sessionUUID := uuid.New()
	params := iamv1.LogoutParams{XSessionUUID: sessionUUID}
	svcErr := errors.New("service error")

	s.sessionService.EXPECT().Logout(s.ctx, sessionUUID).Return(svcErr)

	res, err := s.api.Logout(s.ctx, params)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(res)
}
