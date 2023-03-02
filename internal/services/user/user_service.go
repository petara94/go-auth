package user

import "go-auth/internal/transport/http/api/dto"

type UserRepository interface {
	Create(user dto.User) (*dto.User, error)
	GetByID(id uint64) (*dto.User, error)
	GetByLogin(login string) (*dto.User, error)
	Update(user dto.User) (*dto.User, error)
	DeleteByID(id uint64) error
}

type userService struct {
	userRepository UserRepository
}

func (s *userService) Create(user dto.User) (*dto.User, error) {
	return s.userRepository.Create(user)
}

func (s *userService) Get(id uint64) (*dto.User, error) {
	return s.userRepository.GetByID(id)
}

func (s *userService) Update(user dto.User) (*dto.User, error) {
	return s.userRepository.Update(user)
}

func (s *userService) Delete(id uint64) error {
	return s.userRepository.DeleteByID(id)
}
