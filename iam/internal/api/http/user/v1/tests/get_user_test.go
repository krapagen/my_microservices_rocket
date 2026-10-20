package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

func (s *UserAPISuite) TestGetUser_Success() {
	userUUID := uuid.New()
	sessionUUID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)
	updatedAt := now.Add(time.Minute)
	user := model.User{
		UUID:      userUUID,
		Login:     "validuser",
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}
	params := iamv1.GetUserParams{UserUUID: userUUID, XSessionUUID: sessionUUID}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, nil)
	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(user, nil)

	res, err := s.api.GetUser(s.ctx, params)

	s.NoError(err)
	s.Require().NotNil(res)
	dto, ok := res.(*iamv1.UserDto)
	s.Require().True(ok)
	s.Equal(userUUID, dto.GetUUID())
	s.Equal("validuser", dto.GetLogin())
	s.Equal(now, dto.GetCreatedAt())
	s.Equal(updatedAt, dto.GetUpdatedAt())
}

func (s *UserAPISuite) TestGetUser_NoUpdatedAt() {
	userUUID := uuid.New()
	sessionUUID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)
	user := model.User{
		UUID:      userUUID,
		Login:     "validuser",
		CreatedAt: now,
	}
	params := iamv1.GetUserParams{UserUUID: userUUID, XSessionUUID: sessionUUID}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, nil)
	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(user, nil)

	res, err := s.api.GetUser(s.ctx, params)

	s.NoError(err)
	s.Require().NotNil(res)
	dto, ok := res.(*iamv1.UserDto)
	s.Require().True(ok)
	s.Equal(now, dto.GetUpdatedAt()) // fallback на CreatedAt
}

func (s *UserAPISuite) TestGetUser_Unauthorized() {
	userUUID := uuid.New()
	sessionUUID := uuid.New()
	params := iamv1.GetUserParams{UserUUID: userUUID, XSessionUUID: sessionUUID}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, errs.ErrSessionNotFound)

	res, err := s.api.GetUser(s.ctx, params)

	s.Error(err)
	s.ErrorIs(err, errs.ErrSessionNotFound)
	s.Nil(res)
	s.userService.AssertNotCalled(s.T(), "GetUser")
}

func (s *UserAPISuite) TestGetUser_EmptySession() {
	userUUID := uuid.New()
	params := iamv1.GetUserParams{UserUUID: userUUID, XSessionUUID: uuid.Nil}

	s.sessionService.EXPECT().WhoAmI(s.ctx, uuid.Nil).Return(model.Session{}, model.User{}, errs.ErrEmptySessionID)

	res, err := s.api.GetUser(s.ctx, params)

	s.Error(err)
	s.ErrorIs(err, errs.ErrEmptySessionID)
	s.Nil(res)
	s.userService.AssertNotCalled(s.T(), "GetUser")
}

func (s *UserAPISuite) TestGetUser_UserNotFound() {
	userUUID := uuid.New()
	sessionUUID := uuid.New()
	params := iamv1.GetUserParams{UserUUID: userUUID, XSessionUUID: sessionUUID}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, nil)
	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(model.User{}, errs.ErrUserNotFound)

	res, err := s.api.GetUser(s.ctx, params)

	s.Error(err)
	s.ErrorIs(err, errs.ErrUserNotFound)
	s.Nil(res)
}

func (s *UserAPISuite) TestGetUser_ServiceError() {
	userUUID := uuid.New()
	sessionUUID := uuid.New()
	svcErr := errors.New("service error")
	params := iamv1.GetUserParams{UserUUID: userUUID, XSessionUUID: sessionUUID}

	s.sessionService.EXPECT().WhoAmI(s.ctx, sessionUUID).Return(model.Session{}, model.User{}, nil)
	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(model.User{}, svcErr)

	res, err := s.api.GetUser(s.ctx, params)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(res)
}
