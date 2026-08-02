package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	iamservice "github.com/krapagen/my_microservices_rocket/iam/internal/service/iam"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/iam/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx         context.Context
	userRepo    *mocks.UserRepository
	sessionRepo *mocks.SessionRepository
	service     iamservice.Service
	sessionTTL  time.Duration
	bcryptCost  int
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = mocks.NewUserRepository(s.T())
	s.sessionRepo = mocks.NewSessionRepository(s.T())
	s.sessionTTL = time.Hour
	s.bcryptCost = bcrypt.MinCost
	s.service = iamservice.New(s.userRepo, s.sessionRepo, s.sessionTTL, s.bcryptCost)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
