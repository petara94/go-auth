package api

import (
	"github.com/gofiber/fiber/v2"
	"go-auth/internal/transport/http/api/dto"
)

type UserService interface {
	Create(u dto.User) (*dto.User, error)
	Get(id uint64) (*dto.User, error)
	Update(u dto.User) (*dto.User, error)
	Delete(id uint64) error
}

type UserGroupService interface {
	Create(u dto.UserGroup) (*dto.UserGroup, error)
	Get(id uint64) (*dto.UserGroup, error)
	Update(u dto.UserGroup) (*dto.UserGroup, error)
	Delete(id uint64) error
}

type AuthService interface {
	Login(auth dto.Auth) (*dto.Session, error)
	Logout(session dto.Session) error
}

type Server struct {
	srv *fiber.App

	UserService UserService
	AuthService AuthService
}
