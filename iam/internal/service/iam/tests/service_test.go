package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/service/input"
)

func hashPassword(s *ServiceSuite, password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	s.NoError(err)
	return string(hash)
}

func (s *ServiceSuite) TestRegister_Success() {
	login := "newuser"
	password := "password123"

	s.userRepo.EXPECT().Create(s.ctx, mock.MatchedBy(func(user model.User) bool {
		return user.Login == login &&
			user.UUID != uuid.Nil &&
			user.PasswordHash != "" &&
			!user.CreatedAt.IsZero()
	})).Return(nil)

	userUUID, err := s.service.Register(s.ctx, input.RegisterInput{
		Login:    login,
		Password: password,
	})
	s.NoError(err)
	s.NotEqual(uuid.Nil, userUUID)
}

func (s *ServiceSuite) TestRegister_EmptyLogin() {
	_, err := s.service.Register(s.ctx, input.RegisterInput{
		Login:    "",
		Password: "password123",
	})
	s.ErrorIs(err, errs.ErrInvalidLogin)
}

func (s *ServiceSuite) TestRegister_EmptyPassword() {
	_, err := s.service.Register(s.ctx, input.RegisterInput{
		Login:    "newuser",
		Password: "",
	})
	s.ErrorIs(err, errs.ErrEmptyCredential)
}

func (s *ServiceSuite) TestRegister_WeakPassword() {
	_, err := s.service.Register(s.ctx, input.RegisterInput{
		Login:    "newuser",
		Password: "short",
	})
	s.ErrorIs(err, errs.ErrWeakPassword)
}

func (s *ServiceSuite) TestRegister_UserAlreadyExists() {
	login := "existing"
	password := "password123"

	s.userRepo.EXPECT().Create(s.ctx, mock.Anything).Return(errs.ErrUserAlreadyExists)

	_, err := s.service.Register(s.ctx, input.RegisterInput{
		Login:    login,
		Password: password,
	})
	s.ErrorIs(err, errs.ErrUserAlreadyExists)
}

func (s *ServiceSuite) TestRegister_RepositoryError() {
	login := "newuser"
	password := "password123"
	repoErr := errors.New("repository error")

	s.userRepo.EXPECT().Create(s.ctx, mock.Anything).Return(repoErr)

	_, err := s.service.Register(s.ctx, input.RegisterInput{
		Login:    login,
		Password: password,
	})
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestLogin_Success() {
	login := "validuser"
	password := "password123"
	userUUID := uuid.New()

	hash := hashPassword(s, password)

	s.userRepo.EXPECT().GetByLogin(s.ctx, login).Return(model.User{
		UUID:         userUUID,
		Login:        login,
		PasswordHash: hash,
	}, nil)
	s.sessionRepo.EXPECT().Set(s.ctx, mock.Anything, mock.MatchedBy(func(session model.Session) bool {
		return session.UserUUID == userUUID &&
			session.Login == login &&
			session.UUID != uuid.Nil &&
			!session.ExpiresAt.IsZero()
	}), s.sessionTTL).Return(nil)

	sessionUUID, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    login,
		Password: password,
	})
	s.NoError(err)
	s.NotEqual(uuid.Nil, sessionUUID)
}

func (s *ServiceSuite) TestLogin_InvalidLogin() {
	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    "",
		Password: "password123",
	})
	s.ErrorIs(err, errs.ErrInvalidLogin)
}

func (s *ServiceSuite) TestLogin_EmptyPassword() {
	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    "validuser",
		Password: "",
	})
	s.ErrorIs(err, errs.ErrEmptyCredential)
}

func (s *ServiceSuite) TestLogin_WeakPassword() {
	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    "validuser",
		Password: "short",
	})
	s.ErrorIs(err, errs.ErrWeakPassword)
}

func (s *ServiceSuite) TestLogin_UserNotFound() {
	login := "unknown"
	password := "password123"

	s.userRepo.EXPECT().GetByLogin(s.ctx, login).Return(model.User{}, errs.ErrUserNotFound)

	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    login,
		Password: password,
	})
	s.ErrorIs(err, errs.ErrInvalidCredentials)
}

func (s *ServiceSuite) TestLogin_WrongPassword() {
	login := "validuser"
	password := "password123"
	wrongPassword := "wrongpass123"

	hash := hashPassword(s, password)

	s.userRepo.EXPECT().GetByLogin(s.ctx, login).Return(model.User{
		UUID:         uuid.New(),
		Login:        login,
		PasswordHash: hash,
	}, nil)

	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    login,
		Password: wrongPassword,
	})
	s.ErrorIs(err, errs.ErrInvalidCredentials)
}

func (s *ServiceSuite) TestLogin_SessionSetError() {
	login := "validuser"
	password := "password123"
	setErr := errors.New("set error")

	hash := hashPassword(s, password)

	s.userRepo.EXPECT().GetByLogin(s.ctx, login).Return(model.User{
		UUID:         uuid.New(),
		Login:        login,
		PasswordHash: hash,
	}, nil)
	s.sessionRepo.EXPECT().Set(s.ctx, mock.Anything, mock.Anything, s.sessionTTL).Return(setErr)

	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    login,
		Password: password,
	})
	s.ErrorIs(err, setErr)
}

func (s *ServiceSuite) TestLogout_Success() {
	sessionUUID := uuid.New()

	s.sessionRepo.EXPECT().Delete(s.ctx, sessionUUID.String()).Return(nil)

	err := s.service.Logout(s.ctx, sessionUUID)
	s.NoError(err)
}

func (s *ServiceSuite) TestLogout_EmptySession() {
	err := s.service.Logout(s.ctx, uuid.Nil)
	s.ErrorIs(err, errs.ErrEmptySessionID)
}

func (s *ServiceSuite) TestLogout_DeleteError() {
	sessionUUID := uuid.New()
	deleteErr := errors.New("delete error")

	s.sessionRepo.EXPECT().Delete(s.ctx, sessionUUID.String()).Return(deleteErr)

	err := s.service.Logout(s.ctx, sessionUUID)
	s.ErrorIs(err, deleteErr)
}

func (s *ServiceSuite) TestWhoAmI_Success() {
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()

	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     "validuser",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}
	user := model.User{
		UUID:      userUUID,
		Login:     "validuser",
		CreatedAt: now,
	}

	s.sessionRepo.EXPECT().Get(s.ctx, sessionUUID.String()).Return(session, nil)
	s.userRepo.EXPECT().GetByUUID(s.ctx, userUUID.String()).Return(user, nil)

	resultSession, resultUser, err := s.service.WhoAmI(s.ctx, sessionUUID)
	s.NoError(err)
	s.Equal(session, resultSession)
	s.Equal(user, resultUser)
}

func (s *ServiceSuite) TestWhoAmI_EmptySession() {
	_, _, err := s.service.WhoAmI(s.ctx, uuid.Nil)
	s.ErrorIs(err, errs.ErrEmptySessionID)
}

func (s *ServiceSuite) TestWhoAmI_SessionNotFound() {
	sessionUUID := uuid.New()

	s.sessionRepo.EXPECT().Get(s.ctx, sessionUUID.String()).Return(model.Session{}, errs.ErrSessionNotFound)

	_, _, err := s.service.WhoAmI(s.ctx, sessionUUID)
	s.ErrorIs(err, errs.ErrSessionNotFound)
}

func (s *ServiceSuite) TestWhoAmI_UserNotFound() {
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()

	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     "validuser",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	s.sessionRepo.EXPECT().Get(s.ctx, sessionUUID.String()).Return(session, nil)
	s.userRepo.EXPECT().GetByUUID(s.ctx, userUUID.String()).Return(model.User{}, errs.ErrUserNotFound)

	_, _, err := s.service.WhoAmI(s.ctx, sessionUUID)
	s.ErrorIs(err, errs.ErrUserNotFound)
}

func (s *ServiceSuite) TestWhoAmI_GetSessionError() {
	sessionUUID := uuid.New()
	repoErr := errors.New("session repo error")

	s.sessionRepo.EXPECT().Get(s.ctx, sessionUUID.String()).Return(model.Session{}, repoErr)

	_, _, err := s.service.WhoAmI(s.ctx, sessionUUID)
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestWhoAmI_GetUserError() {
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()
	repoErr := errors.New("user repo error")

	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		Login:     "validuser",
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	s.sessionRepo.EXPECT().Get(s.ctx, sessionUUID.String()).Return(session, nil)
	s.userRepo.EXPECT().GetByUUID(s.ctx, userUUID.String()).Return(model.User{}, repoErr)

	_, _, err := s.service.WhoAmI(s.ctx, sessionUUID)
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestLogin_GetUserError() {
	login := "validuser"
	password := "password123"
	repoErr := errors.New("repo error")

	s.userRepo.EXPECT().GetByLogin(s.ctx, login).Return(model.User{}, repoErr)

	_, err := s.service.Login(s.ctx, input.LoginInput{
		Login:    login,
		Password: password,
	})
	s.ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestGetUser_Success() {
	userUUID := uuid.New()
	user := model.User{
		UUID:      userUUID,
		Login:     "validuser",
		CreatedAt: time.Now(),
	}

	s.userRepo.EXPECT().GetByUUID(s.ctx, userUUID.String()).Return(user, nil)

	result, err := s.service.GetUser(s.ctx, userUUID)
	s.NoError(err)
	s.Equal(user, result)
}

func (s *ServiceSuite) TestGetUser_EmptyUUID() {
	_, err := s.service.GetUser(s.ctx, uuid.Nil)
	s.ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestGetUser_UserNotFound() {
	userUUID := uuid.New()

	s.userRepo.EXPECT().GetByUUID(s.ctx, userUUID.String()).Return(model.User{}, errs.ErrUserNotFound)

	_, err := s.service.GetUser(s.ctx, userUUID)
	s.ErrorIs(err, errs.ErrUserNotFound)
}
