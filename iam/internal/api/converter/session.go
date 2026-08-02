package converter

import (
	"errors"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
	commonv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/common/v1"
)

func SessionToDto(session model.Session) (*commonv1.Session, error) {
	if session.UUID == uuid.Nil {
		return nil, errors.New("SessionToDto: session.UUID равен nil")
	}
	dto := &commonv1.Session{
		Uuid:      session.UUID.String(),
		CreatedAt: timestamppb.New(session.CreatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
	if session.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*session.UpdatedAt)
	}
	return dto, nil
}

func LoginRequestToInput(req *authv1.LoginRequest) (input.LoginInput, error) {
	login := req.GetLogin()
	password := req.GetPassword()
	if login == "" {
		return input.LoginInput{}, errs.ErrInvalidLogin
	}
	if password == "" {
		return input.LoginInput{}, errs.ErrEmptyCredential
	}
	if len(password) < 8 {
		return input.LoginInput{}, errs.ErrWeakPassword
	}

	return input.LoginInput{
		Login:    login,
		Password: password,
	}, nil
}

func LogoutRequestToInput(req *authv1.LogoutRequest) (uuid.UUID, error) {
	if req.GetSessionUuid() == "" {
		return uuid.Nil, errs.ErrEmptySessionID
	}
	sessionUUID, err := uuid.Parse(req.GetSessionUuid())
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return sessionUUID, nil
}

func WhoAmIRequestToInput(req *authv1.WhoamiRequest) (uuid.UUID, error) {
	if req.GetSessionUuid() == "" {
		return uuid.Nil, errs.ErrEmptySessionID
	}
	sessionUUID, err := uuid.Parse(req.GetSessionUuid())
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return sessionUUID, nil
}
