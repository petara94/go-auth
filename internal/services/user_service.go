package services

import (
	"context"
	"github.com/petara94/go-auth/internal/services/dto"
	"github.com/petara94/go-auth/internal/services/pkg"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2 --name UserRepository
type UserRepository interface {
	Create(user dto.User) (uint64, error)
	GetByID(id uint64) (*dto.User, error)
	GetWithPagination(perPage, page int) ([]*dto.User, error)
	GetByLogin(login string) (*dto.User, error)
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
	// check password is strong
	if !pkg.CheckPassword(user.Password) && user.CheckPassword {
		return nil, ErrWeakPassword
	}

	user.Password = pkg.HashPassword(user.Password)

	id, err := s.userRepository.Create(user)
	if err != nil {
		return nil, HandleRepositoryError(err)
	}

	return s.userRepository.GetByID(id)
}

// CreateWithLoginAndPassword is a helper function for creating user with login and password
func (s *UserService) CreateWithLoginAndPassword(login, password string) (*dto.User, error) {
	return s.Create(dto.User{
		Login:         login,
		Password:      password,
		CheckPassword: true,
	})
}

// GetByID returns user by id
func (s *UserService) GetByID(id uint64) (*dto.User, error) {
	return s.userRepository.GetByID(id)
}

// GetByLogin returns user by login
func (s *UserService) GetByLogin(login string) (*dto.User, error) {
	return s.userRepository.GetByLogin(login)
}

// Update updates user
func (s *UserService) Update(user dto.User) (*dto.User, error) {
	old, err := s.userRepository.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	// check is password is changed
	if pkg.PasswordEqual(user.Password, old.Password) {
		user.Password = old.Password
	} else {
		user.Password = pkg.HashPassword(user.Password)
	}

	return s.userRepository.Update(user)
}

// UpdatePassword updates user password
func (s *UserService) UpdatePassword(id uint64, oldPassword, newPassword string) error {
	if oldPassword == newPassword {
		return ErrSamePassword
	}

	// check if password is strong
	if !pkg.CheckPassword(newPassword) {
		return ErrWeakPassword
	}

	// get user
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	// check if old password is correct
	if !pkg.PasswordEqual(oldPassword, user.Password) {
		return ErrWrongPassword
	}

	// update password
	user.Password = pkg.HashPassword(newPassword)
	_, err = s.userRepository.Update(*user)
	if err != nil {
		return err
	}

	return nil
}

// CreateWithLogin is a helper function for creating user with login and empty password
func (s *UserService) CreateWithLogin(login string) (*dto.User, error) {
	// check if user with this login already exists
	_, err := s.userRepository.GetByLogin(login)
	if err == nil {
		return nil, ErrLoginAlreadyExists
	}

	// create user
	user := dto.User{
		Login: login,
	}

	id, err := s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return s.userRepository.GetByID(id)
}

// SetCheckPassword sets user.UsePasswordConstraints to usePasswordConstraints
func (s *UserService) SetCheckPassword(id uint64, checkPassword bool) error {
	// get user
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	// check if user.CheckPassword is already equal to checkPassword
	if user.CheckPassword == checkPassword {
		return nil
	}

	// update user
	user.CheckPassword = checkPassword
	_, err = s.userRepository.Update(*user)

	return err
}

// GetWithPagination returns users with pagination
func (s *UserService) GetWithPagination(perPage, page int) ([]*dto.User, error) {
	// check pagination
	if perPage < 1 || page < 1 {
		return nil, ErrWrongPagination
	}

	return s.userRepository.GetWithPagination(perPage, page)
}

// SetAdmin sets user.IsAdmin to isAdnmin
func (s *UserService) SetAdmin(id uint64, isAdnmin bool) error {
	// check if user exists
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	// check if user.IsAdmin is already equal to isAdnmin
	if user.IsAdmin == isAdnmin {
		return nil
	}

	// update user
	user.IsAdmin = isAdnmin
	_, err = s.userRepository.Update(*user)

	return err
}

// SetBlockUser sets user.IsBlocked to isBlock
func (s *UserService) SetBlockUser(id uint64, isBlock bool) error {
	// check if user exists
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	// check if user.IsBlocked is already equal to isBlock
	if user.IsBlocked == isBlock {
		return nil
	}

	// update user
	user.IsBlocked = isBlock
	_, err = s.userRepository.Update(*user)

	return err
}

// Delete deletes user by id
func (s *UserService) Delete(id uint64) error {
	return s.userRepository.DeleteByID(id)
}
