package interceptor

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
)

func ErrorInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("перехвачена паника",
				"panic", r,
				"method", info.FullMethod,
				"stack", string(debug.Stack()))
			err = status.Error(codes.Internal, "внутренняя ошибка")
		}
	}()

	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	switch {
	case errors.Is(err, errs.ErrInvalidLogin),
		errors.Is(err, errs.ErrWeakPassword),
		errors.Is(err, errs.ErrEmptyCredential),
		errors.Is(err, errs.ErrEmptySessionID),
		errors.Is(err, errs.ErrInvalidUUID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrUserAlreadyExists):
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrSessionNotFound):
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, errs.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
