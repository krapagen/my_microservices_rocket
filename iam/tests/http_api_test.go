//go:build apitest

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// uniqueLogin генерирует уникальный логин, чтобы кейсы не пересекались в общей БД
func uniqueLogin(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, uuid.NewString()[:8])
}

// registerHTTP регистрирует пользователя и возвращает user_uuid
func registerHTTP(t *testing.T, login, password string) uuid.UUID {
	t.Helper()
	res, err := httpClient.Register(context.Background(), &iamv1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)

	created, ok := res.(*iamv1.RegisterResponse)
	require.True(t, ok, "ожидался RegisterResponse, получили %T", res)
	require.NotEqual(t, uuid.Nil, created.GetUserUUID())

	return created.GetUserUUID()
}

// loginHTTP логинит пользователя и возвращает session_uuid
func loginHTTP(t *testing.T, login, password string) uuid.UUID {
	t.Helper()
	res, err := httpClient.Login(context.Background(), &iamv1.LoginRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)

	session, ok := res.(*iamv1.LoginResponse)
	require.True(t, ok, "ожидался LoginResponse, получили %T", res)
	require.NotEqual(t, uuid.Nil, session.GetSessionUUID())

	return session.GetSessionUUID()
}

func TestHTTP_RegisterLoginGetLogout_HappyPath(t *testing.T) {
	ctx := context.Background()
	login := uniqueLogin("http-alice")

	userUUID := registerHTTP(t, login, "password123")
	sessionUUID := loginHTTP(t, login, "password123")

	// GetUser с валидной сессией — 200
	userRes, err := httpClient.GetUser(ctx, iamv1.GetUserParams{
		UserUUID:     userUUID,
		XSessionUUID: sessionUUID,
	})
	require.NoError(t, err)

	user, ok := userRes.(*iamv1.UserDto)
	require.True(t, ok, "ожидался UserDto, получили %T", userRes)
	require.Equal(t, userUUID, user.GetUUID())
	require.Equal(t, login, user.GetLogin())

	// Logout — 204
	logoutRes, err := httpClient.Logout(ctx, iamv1.LogoutParams{
		XSessionUUID: sessionUUID,
	})
	require.NoError(t, err)
	_, ok = logoutRes.(*iamv1.LogoutNoContent)
	require.True(t, ok, "ожидался LogoutNoContent, получили %T", logoutRes)

	// GetUser после logout — 401
	userRes, err = httpClient.GetUser(ctx, iamv1.GetUserParams{
		UserUUID:     userUUID,
		XSessionUUID: sessionUUID,
	})
	require.NoError(t, err)
	_, ok = userRes.(*iamv1.GetUserUnauthorized)
	require.True(t, ok, "ожидался GetUserUnauthorized, получили %T", userRes)
}

func TestHTTP_Register_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	login := uniqueLogin("dup")

	registerHTTP(t, login, "password123")

	res, err := httpClient.Register(ctx, &iamv1.RegisterRequest{
		Login:    login,
		Password: "password123",
	})
	require.NoError(t, err)

	conflict, ok := res.(*iamv1.RegisterConflict)
	require.True(t, ok, "ожидался RegisterConflict, получили %T", res)
	require.NotEmpty(t, conflict.Message)
}

func TestHTTP_Register_ValidationError(t *testing.T) {
	ctx := context.Background()

	res, err := httpClient.Register(ctx, &iamv1.RegisterRequest{
		Login:    uniqueLogin("weak"),
		Password: "short",
	})
	require.NoError(t, err)

	bad, ok := res.(*iamv1.RegisterBadRequest)
	require.True(t, ok, "ожидался RegisterBadRequest, получили %T", res)
	require.NotEmpty(t, bad.Message)
}

func TestHTTP_Login_InvalidCredentials(t *testing.T) {
	ctx := context.Background()

	res, err := httpClient.Login(ctx, &iamv1.LoginRequest{
		Login:    "ghost-" + uuid.NewString()[:8],
		Password: "any-password",
	})
	require.NoError(t, err)

	unauth, ok := res.(*iamv1.LoginUnauthorized)
	require.True(t, ok, "ожидался LoginUnauthorized, получили %T", res)
	require.NotEmpty(t, unauth.Message)
}

func TestHTTP_GetUser_WithoutSession(t *testing.T) {
	ctx := context.Background()

	res, err := httpClient.GetUser(ctx, iamv1.GetUserParams{
		UserUUID:     uuid.New(),
		XSessionUUID: uuid.New(), // несуществующая сессия
	})
	require.NoError(t, err)

	unauth, ok := res.(*iamv1.GetUserUnauthorized)
	require.True(t, ok, "ожидался GetUserUnauthorized, получили %T", res)
	require.NotEmpty(t, unauth.Message)
}

func TestHTTP_GetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	login := uniqueLogin("lookup")

	registerHTTP(t, login, "password123")
	sessionUUID := loginHTTP(t, login, "password123")

	res, err := httpClient.GetUser(ctx, iamv1.GetUserParams{
		UserUUID:     uuid.New(), // несуществующий пользователь
		XSessionUUID: sessionUUID,
	})
	require.NoError(t, err)

	notFound, ok := res.(*iamv1.GetUserNotFound)
	require.True(t, ok, "ожидался GetUserNotFound, получили %T", res)
	require.NotEmpty(t, notFound.Message)
}

func TestHTTP_Logout_Idempotent(t *testing.T) {
	ctx := context.Background()

	// Logout несуществующей сессии не должен возвращать ошибку
	res, err := httpClient.Logout(ctx, iamv1.LogoutParams{
		XSessionUUID: uuid.New(),
	})
	require.NoError(t, err)

	_, ok := res.(*iamv1.LogoutNoContent)
	require.True(t, ok, "ожидался LogoutNoContent, получили %T", res)
}
