package services

import (
	"context"
	"fmt"
	"github.com/petara94/go-auth/internal/services/pkg"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"go.uber.org/zap"
)

//go:generate mockery --name UserRepository

type UserRepository interface {
	Create(user dto.User) (uint64, error)
	GetByID(id uint64) (*dto.User, error)
	Get(perPage, page int) ([]*dto.User, error)
	GetByLogin(login string) (*dto.User, error)
	LinkToGroup(id, groupId uint64) error
	Update(user dto.User) (*dto.User, error)
	DeleteByID(id uint64) error
}

type UserService struct {
	ctx            context.Context
	userRepository UserRepository
	logger         zap.Logger
}

func NewUserService(ctx context.Context, userRepository UserRepository, logger zap.Logger) *UserService {
	return &UserService{userRepository: userRepository, logger: logger, ctx: ctx}
}

func (s *UserService) Create(user dto.User) (*dto.User, error) {
	user.Password = pkg.HashPassword(user.Password)
	id, err := s.userRepository.Create(user)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	if user.UserGroupID != nil {
		err = s.userRepository.LinkToGroup(id, *user.UserGroupID)
		if err != nil {
			return nil, fmt.Errorf("linking user to user group: %w", err)
		}
	}

	return s.userRepository.GetByID(id)
}

func (s *UserService) GetByID(id uint64) (*dto.User, error) {
	return s.userRepository.GetByID(id)
}

func (s *UserService) Get(perPage, page int) ([]*dto.User, error) {
	return s.userRepository.Get(perPage, page)
}

func (s *UserService) GetByLogin(login string) (*dto.User, error) {
	return s.userRepository.GetByLogin(login)
}

func (s *UserService) Update(user dto.User) (*dto.User, error) {
	old, err := s.userRepository.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	if pkg.PasswordEqual(user.Password, old.Password) {
		user.Password = old.Password
	} else {
		user.Password = pkg.HashPassword(user.Password)
	}

	return s.userRepository.Update(user)
}

func (s *UserService) Delete(id uint64) error {
	return s.userRepository.DeleteByID(id)
}
