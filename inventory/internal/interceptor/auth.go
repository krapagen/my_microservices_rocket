package interceptor

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	iamv1 "github.com/krapagen/my_microservices_rocket/inventory/internal/client/grpc/iam/v1"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/auth"
)

// SessionMetadataKey — ключ, по которому interceptor передаёт session-uuid в gRPC metadata.
// Экспортирован, чтобы его могли использовать тесты и другие потребители пакета.
const SessionMetadataKey = "session-uuid"

// publicMethods — методы gRPC Reflection API, не требующие аутентификации
// Reflection используется инструментами вроде grpcurl и grpcui для интроспекции сервера
// (получение списка сервисов, описания методов и т.д.)
// Указаны обе версии (v1 и v1alpha), т.к. разные клиенты используют разные версии API
var publicMethods = map[string]bool{
	"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      true,
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
}

// GRPCAuth — gRPC unary server interceptor для аутентификации
// Извлекает session-token из incoming metadata, проверяет сессию
// и добавляет пользователя в контекст
//
// gRPC metadata — аналог HTTP-заголовков, бывает двух видов:
//   - Incoming metadata — то, что пришло от клиента (на стороне сервера)
//     Читается через metadata.FromIncomingContext(ctx)
//     Используется для чтения токенов, заголовков авторизации, request-id и т.д
//   - Outgoing metadata — то, что клиент отправляет серверу (на стороне клиента)
//     Устанавливается через metadata.NewOutgoingContext(ctx, md)
//     Используется для передачи токенов, трейс-контекста и прочих заголовков при вызове
//
// Разделение нужно, чтобы сервер случайно не переслал чужие заголовки дальше по цепочке
// вызовов — incoming и outgoing живут в разных «слотах» контекста
//
// Пример: Client → ServiceA → ServiceB
// ServiceA получает request-id через incoming metadata от клиента
// Чтобы пробросить его дальше в ServiceB, нужно явно переложить значение
// из incoming в outgoing:
//
//	md, _ := metadata.FromIncomingContext(ctx)
//	outCtx := metadata.NewOutgoingContext(ctx, md)
//	serviceB.Call(outCtx, req)
//
// Если этого не сделать, ServiceB не увидит request-id — incoming автоматически
// не попадает в outgoing. Это сделано намеренно, чтобы не пробрасывать
// всё подряд (например, auth-токен клиента туда, куда он не должен попасть)
// Auth — alias for GRPCAuth. Keeps the interceptor API surface short for tests and consumers that
// don't need the transport-specific prefix.
func Auth(iamClient iamv1.Client) grpc.UnaryServerInterceptor {
	return GRPCAuth(iamClient)
}

func GRPCAuth(iamClient iamv1.Client) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Пропускаем аутентификацию для публичных методов
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "отсутствуют metadata")
		}

		sessionUUIDs := md.Get(SessionMetadataKey)
		if len(sessionUUIDs) == 0 {
			return nil, status.Error(codes.Unauthenticated, "отсутствует session-uuid в metadata")
		}

		sessionUUID := sessionUUIDs[0]
		if sessionUUID == "" {
			return nil, status.Error(codes.Unauthenticated, "пустой session-uuid")
		}

		user, _, err := iamClient.Whoami(ctx, sessionUUID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "недействительный session UUID")
		}

		userUuid, err := uuid.Parse(user.Uuid)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "не удалось распарсить user UUID")
		}
		ctx = auth.WithUserUUID(ctx, userUuid)
		slog.Info(
			"gRPC: пользователь аутентифицирован",
			"method", info.FullMethod,
			"user_id", user.Uuid,
		)
		return handler(ctx, req)
	}
}
