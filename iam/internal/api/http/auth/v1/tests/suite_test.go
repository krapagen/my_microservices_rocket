package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	authv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/auth/v1"
	"github.com/krapagen/my_microservices_rocket/iam/internal/api/http/auth/v1/mocks"
)

type AuthAPISuite struct {
	suite.Suite
	ctx            context.Context
	sessionService *mocks.SessionService
	api            authv1api.AuthHandler
}

func (s *AuthAPISuite) SetupTest() {
	s.ctx = context.Background()
	s.sessionService = mocks.NewSessionService(s.T())
	s.api = authv1api.New(s.sessionService)
}

func TestAuthAPISuite(t *testing.T) {
	suite.Run(t, new(AuthAPISuite))
}
