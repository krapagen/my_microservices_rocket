package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
)

func (s *AuthAPISuite) TestWhoami_Success() {
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()
	updatedAt := now.Add(time.Minute)

	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     "validuser",
		CreatedAt: now,
		UpdatedAt: &updatedAt,
		ExpiresAt: now.Add(time.Hour),
	}
	user := model.User{
		UUID:      userUUID,
		Login:     "validuser",
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}
	req := &authv1.WhoamiRequest{SessionUuid: sessionUUID.String()}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(session, user, nil)

	resp, err := s.api.Whoami(s.ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(sessionUUID.String(), resp.GetSession().GetUuid())
	s.Equal(userUUID.String(), resp.GetUser().GetUuid())
}

func (s *AuthAPISuite) TestWhoami_EmptySession() {
	req := &authv1.WhoamiRequest{SessionUuid: ""}

	resp, err := s.api.Whoami(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptySessionID)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestWhoami_InvalidUUID() {
	req := &authv1.WhoamiRequest{SessionUuid: "not-a-uuid"}

	resp, err := s.api.Whoami(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestWhoami_SessionNotFound() {
	sessionUUID := uuid.New()
	req := &authv1.WhoamiRequest{SessionUuid: sessionUUID.String()}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, errs.ErrSessionNotFound)

	resp, err := s.api.Whoami(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrSessionNotFound)
	s.Nil(resp)
}

func (s *AuthAPISuite) TestWhoami_ServiceError() {
	sessionUUID := uuid.New()
	req := &authv1.WhoamiRequest{SessionUuid: sessionUUID.String()}
	svcErr := errors.New("service error")

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, svcErr)

	resp, err := s.api.Whoami(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}
