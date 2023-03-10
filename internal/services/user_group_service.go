package services

import (
	"context"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
)

//go:generate mockery --name UserGroupRepository
type UserGroupRepository interface {
	Create(group dto.UserGroup) (*dto.UserGroup, error)
	GetByID(id uint64) (*dto.UserGroup, error)
	Update(group dto.UserGroup) (*dto.UserGroup, error)
	DeleteByID(id uint64) error
}

type UserGroupService struct {
	ctx                 context.Context
	userGroupRepository UserGroupRepository
}

func (s *UserGroupService) Create(group dto.UserGroup) (*dto.UserGroup, error) {
	return s.userGroupRepository.Create(group)
}

func (s *UserGroupService) Get(id uint64) (*dto.UserGroup, error) {
	return s.userGroupRepository.GetByID(id)
}

func (s *UserGroupService) Update(group dto.UserGroup) (*dto.UserGroup, error) {
	return s.userGroupRepository.Update(group)
}

func (s *UserGroupService) Delete(id uint64) error {
	return s.userGroupRepository.DeleteByID(id)
}
