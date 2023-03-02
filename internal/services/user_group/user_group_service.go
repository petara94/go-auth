package user_group

import "go-auth/internal/transport/http/api/dto"

type UserGroupRepository interface {
	Create(group dto.UserGroup) (*dto.UserGroup, error)
	GetByID(id uint64) (*dto.UserGroup, error)
	Update(group dto.UserGroup) (*dto.UserGroup, error)
	DeleteByID(id uint64) error
}

type userGroupService struct {
	userGroupRepository UserGroupRepository
}

func (s *userGroupService) Create(group dto.UserGroup) (*dto.UserGroup, error) {
	return s.userGroupRepository.Create(group)
}

func (s *userGroupService) Get(id uint64) (*dto.UserGroup, error) {
	return s.userGroupRepository.GetByID(id)
}

func (s *userGroupService) Update(group dto.UserGroup) (*dto.UserGroup, error) {
	return s.userGroupRepository.Update(group)
}

func (s *userGroupService) Delete(id uint64) error {
	return s.userGroupRepository.DeleteByID(id)
}
