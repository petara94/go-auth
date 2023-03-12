package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/petara94/go-auth/internal/services/pkg"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"time"
)

//go:generate mockery --name SessionRepository
type SessionRepository interface {
	Create(session dto.Session) (*dto.Session, error)
	GetByToken(token string) (*dto.Session, error)
	DeleteByToken(token string) error
}

type AuthService struct {
	ctx                 context.Context
	sessionRepository   SessionRepository
	userGroupRepository UserRepository
}

func NewAuthService(ctx context.Context, sessionRepository SessionRepository, userGroupRepository UserRepository) *AuthService {
	return &AuthService{ctx: ctx, sessionRepository: sessionRepository, userGroupRepository: userGroupRepository}
}

func (s *AuthService) Login(auth dto.Auth) (*dto.Session, error) {
	userByLogin, err := s.userGroupRepository.GetByLogin(auth.Login)
	if err != nil {
		return nil, err
	}

	if !pkg.PasswordEqual(auth.Password, userByLogin.Password) {
		return nil, ErrLoginErr
	}

	session, err := s.sessionRepository.Create(dto.Session{
		Token:  uuid.NewString(),
		UserID: userByLogin.ID,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Logout(session dto.Session) error {
	return s.sessionRepository.DeleteByToken(session.Token)
}

func (s *AuthService) Get(token string) (*dto.Session, error) {
	session, err := s.sessionRepository.GetByToken(token)
	if err != nil {
		return nil, err
	}

	if session.Expr != nil {
		if session.Expr.Unix() > time.Now().Unix() {
			return nil, ErrSessionExpired
		}
	}

	return session, nil
}
