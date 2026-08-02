package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/platform/pkg/auth"
)

type iamClient interface {
	Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error)
}

// HTTPAuth — HTTP middleware, который извлекает Bearer-токен из заголовка Authorization,
// проверяет сессию через IAM и кладёт user UUID / session UUID в контекст запроса.
func HTTPAuth(authClient iamClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "отсутствует заголовок Authorization", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "формат: Bearer <session-uuid>", http.StatusUnauthorized)
				return
			}

			sessionUUID := strings.TrimPrefix(authHeader, "Bearer ")
			if sessionUUID == "" {
				http.Error(w, "пустой session-uuid", http.StatusUnauthorized)
				return
			}

			userUUID, err := authClient.Whoami(r.Context(), sessionUUID)
			if err != nil {
				http.Error(w, "недействительная сессия", http.StatusUnauthorized)
				return
			}

			ctx := auth.WithSessionUUID(r.Context(), sessionUUID)
			ctx = auth.WithUserUUID(ctx, userUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
