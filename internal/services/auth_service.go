package services

import (
	"context"
	"github.com/petara94/go-auth/internal/services/dto"
	"time"

	"github.com/google/uuid"
	"github.com/petara94/go-auth/internal/services/pkg"
)

//go:generate mockery --name SessionRepository
type SessionRepository interface {
	GetWithPagination(perPage, page int) ([]*dto.Session, error)
	Create(session dto.Session) (*dto.Session, error)
	GetByToken(token string) (*dto.Session, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uint64) error
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

	// if blocked
	if userByLogin.IsBlocked {
		return nil, ErrUserBlocked
	}

	if !pkg.PasswordEqual(auth.Password, userByLogin.Password) {
		return nil, ErrLoginErr
	}

	session, err := s.sessionRepository.Create(dto.Session{
		Token:  uuid.NewString(),
		UserID: userByLogin.ID,
		Expr:   auth.TTL,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Logout(session dto.Session) error {
	return s.sessionRepository.DeleteByToken(session.Token)
}

func (s *AuthService) GetWithPagination(perPage, page int) ([]*dto.Session, error) {
	// check pagination
	if perPage < 1 || page < 1 {
		return nil, ErrWrongPagination
	}

	return s.sessionRepository.GetWithPagination(perPage, page)
}

func (s *AuthService) Get(token string) (*dto.Session, error) {
	session, err := s.sessionRepository.GetByToken(token)
	if err != nil {
		return nil, err
	}

	if session.Expr != nil {
		if session.Expr.Unix() < time.Now().Unix() {
			return nil, ErrSessionExpired
		}
	}

	return session, nil
}

func (s *AuthService) DeleteByUserID(userID uint64) error {
	return s.sessionRepository.DeleteByUserID(userID)
}
