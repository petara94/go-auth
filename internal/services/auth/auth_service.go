package auth

import (
	"github.com/google/uuid"
	"go-auth/internal/services/pkg"
	"go-auth/internal/services/user"
	"go-auth/internal/transport/http/api/dto"
)

type SessionRepository interface {
	Create(session dto.Session) (*dto.Session, error)
	GetByToken(token string) (*dto.Session, error)
	DeleteByToken(token string) error
}

type authService struct {
	sessionRepository   SessionRepository
	userGroupRepository user.UserRepository
}

func (s *authService) Login(auth dto.Auth) (*dto.Session, error) {
	byLogin, err := s.userGroupRepository.GetByLogin(auth.Login)
	if err != nil {
		return nil, err
	}

	if pkg.PasswordEqual(auth.Password, byLogin.Password) {
		return nil, err
	}

	session, err := s.sessionRepository.Create(dto.Session{
		Token:  uuid.NewString(),
		UserID: byLogin.ID,
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *authService) Logout(session dto.Session) error {
	return s.sessionRepository.DeleteByToken(session.Token)
}
