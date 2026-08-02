package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/krapagen/my_microservices_rocket/platform/pkg/auth"
)

// SessionForwarder — исходящий (client-side) gRPC-интерцептор,
// который достаёт session_uuid из контекста и прокидывает его
// в исходящие gRPC-вызовы через metadata.
func SessionForwarder() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		sessionUUID, ok := auth.SessionUUIDFromContext(ctx)
		if ok && sessionUUID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "session-uuid", sessionUUID)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
