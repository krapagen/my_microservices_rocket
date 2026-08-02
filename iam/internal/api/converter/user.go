package converter

import (
	"errors"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	commonv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/common/v1"
	userv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

func UserToDto(user model.User) (*commonv1.User, error) {
	if user.UUID == uuid.Nil {
		return nil, errors.New("UserToDto: user.UUID равен nil")
	}
	dto := &commonv1.User{
		Uuid: user.UUID.String(),
		Info: &commonv1.UserInfo{
			Login: user.Login,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	if user.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}
	return dto, nil
}

func RegisterRequestToInput(req *userv1.RegisterRequest) (input.RegisterInput, error) {
	login := req.GetInfo().GetInfo().GetLogin()
	password := req.GetInfo().GetPassword()
	if login == "" {
		return input.RegisterInput{}, errs.ErrInvalidLogin
	}
	if password == "" {
		return input.RegisterInput{}, errs.ErrEmptyCredential
	}
	if len(password) < 8 {
		return input.RegisterInput{}, errs.ErrWeakPassword
	}

	return input.RegisterInput{
		Login:    login,
		Password: password,
	}, nil
}

func GetUserRequestToInput(req *userv1.GetUserRequest) (uuid.UUID, error) {
	if req.GetUserUuid() == "" {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return userUUID, nil
}
