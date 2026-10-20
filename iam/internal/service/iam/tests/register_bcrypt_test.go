package tests

import (
	"golang.org/x/crypto/bcrypt"

	iamservice "github.com/krapagen/my_microservices_rocket/iam/internal/service/iam"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

// TestRegister_BcryptError — bcrypt.GenerateFromPassword падает при cost выше
// допустимого MaxCost, и сервис должен вернуть ошибку хэширования
func (s *ServiceSuite) TestRegister_BcryptError() {
	svc := iamservice.New(s.userRepo, s.sessionRepo, s.sessionTTL, bcrypt.MaxCost+1)

	_, err := svc.Register(s.ctx, input.RegisterInput{
		Login:    "newuser",
		Password: "password123",
	})

	s.Error(err)
	s.userRepo.AssertNotCalled(s.T(), "Create")
}
