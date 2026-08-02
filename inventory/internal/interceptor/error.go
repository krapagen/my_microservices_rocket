package interceptor

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

// internal/interceptor/error.go
func ErrorInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("перехвачена паника",
				"panic", r,
				"method", info.FullMethod,
				"stack", string(debug.Stack()))
		}
	}()
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	if _, ok := status.FromError(err); ok {
		return nil, err
	}
	switch {
	case errors.Is(err, errs.ErrPartNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, errs.ErrInvalidUUID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrPartTypeMismatch):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrIncompatibleParts):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrOutOfStock):
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, errs.ErrNothingToRelease):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrCommitParts):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrInvalidProperties):
		return nil, status.Error(codes.Internal, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
