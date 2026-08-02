package iam

import "time"

type service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	sessionTTL  time.Duration
	bcryptCost  int
}

// New создаёт сервис управления пользователями и сессиями
func New(userRepo UserRepository, sessionRepo SessionRepository, sessionTTL time.Duration, bcryptCost int) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
		bcryptCost:  bcryptCost,
	}
}
