package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

func (s *UserAPISuite) TestGetUser_Success() {
	userUUID := uuid.New()
	now := time.Now()
	updatedAt := now.Add(time.Minute)
	user := model.User{
		UUID:      userUUID,
		Login:     "validuser",
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}
	req := &userv1.GetUserRequest{UserUuid: userUUID.String()}

	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(user, nil)

	resp, err := s.api.GetUser(s.ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(userUUID.String(), resp.GetUser().GetUuid())
	s.Equal("validuser", resp.GetUser().GetInfo().GetLogin())
}

func (s *UserAPISuite) TestGetUser_EmptyUUID() {
	req := &userv1.GetUserRequest{UserUuid: ""}

	resp, err := s.api.GetUser(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *UserAPISuite) TestGetUser_InvalidUUID() {
	req := &userv1.GetUserRequest{UserUuid: "not-a-uuid"}

	resp, err := s.api.GetUser(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, errs.ErrInvalidUUID)
	s.Nil(resp)
}

func (s *UserAPISuite) TestGetUser_ServiceError() {
	userUUID := uuid.New()
	svcErr := errors.New("service error")
	req := &userv1.GetUserRequest{UserUuid: userUUID.String()}

	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(model.User{}, svcErr)

	resp, err := s.api.GetUser(s.ctx, req)

	s.Error(err)
	s.ErrorIs(err, svcErr)
	s.Nil(resp)
}

func (s *UserAPISuite) TestGetUser_ConvertError() {
	userUUID := uuid.New()
	req := &userv1.GetUserRequest{UserUuid: userUUID.String()}

	s.userService.EXPECT().GetUser(s.ctx, userUUID).Return(model.User{}, nil)

	resp, err := s.api.GetUser(s.ctx, req)

	s.Error(err)
	s.Nil(resp)
}
