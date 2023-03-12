package services

import (
	"context"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"go.uber.org/zap"
)

//go:generate mockery --name UserGroupRepository
type UserGroupRepository interface {
	Create(group dto.UserGroup) (*dto.UserGroup, error)
	GetByID(id uint64) (*dto.UserGroup, error)
	Get(perPage, page int) ([]*dto.UserGroup, error)
	Update(group dto.UserGroup) (*dto.UserGroup, error)
	DeleteByID(id uint64) error
}

type UserGroupService struct {
	ctx                 context.Context
	userGroupRepository UserGroupRepository
	logger              zap.Logger
}

func NewUserGroupService(ctx context.Context, userGroupRepository UserGroupRepository, logger zap.Logger) *UserGroupService {
	return &UserGroupService{ctx: ctx, userGroupRepository: userGroupRepository, logger: logger}
}

func (s *UserGroupService) Create(group dto.UserGroup) (*dto.UserGroup, error) {
	return s.userGroupRepository.Create(group)
}

func (s *UserGroupService) GetByID(id uint64) (*dto.UserGroup, error) {
	return s.userGroupRepository.GetByID(id)
}

func (s *UserGroupService) Get(perPage, page int) ([]*dto.UserGroup, error) {
	return s.userGroupRepository.Get(perPage, page)
}

func (s *UserGroupService) Update(group dto.UserGroup) (*dto.UserGroup, error) {
	return s.userGroupRepository.Update(group)
}

func (s *UserGroupService) Delete(id uint64) error {
	return s.userGroupRepository.DeleteByID(id)
}
