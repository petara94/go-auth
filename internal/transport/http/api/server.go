package api

import (
	"github.com/gofiber/fiber/v2"
	serv_dto "github.com/petara94/go-auth/internal/services/dto"
	"go.uber.org/zap"
)

type UserService interface {
	Create(u serv_dto.User) (*serv_dto.User, error)
	CreateWithLoginAndPassword(login, password string) (*serv_dto.User, error)
	GetByID(id uint64) (*serv_dto.User, error)
	GetByLogin(login string) (*serv_dto.User, error)
	UpdatePassword(id uint64, oldPassword, newPassword string) error
	Delete(id uint64) error

	// Admin access
	CreateWithLogin(login string) (*serv_dto.User, error)
	GetWithPagination(perPage, page int) ([]*serv_dto.User, error)
	SetAdmin(id uint64, isAdnmin bool) error
	SetCheckPassword(id uint64, usePasswordConstraint bool) error
	SetBlockUser(id uint64, isBlock bool) error
	Update(u serv_dto.User) (*serv_dto.User, error)
}

type AuthService interface {
	Login(auth serv_dto.Auth) (*serv_dto.Session, error)
	Get(token string) (*serv_dto.Session, error)
	Logout(session serv_dto.Session) error
	DeleteByUserID(userID uint64) error
}

type Server struct {
	srv  *fiber.App
	conf *ServerConfig

	logger zap.Logger

	UserService UserService
	AuthService AuthService
}

type ServerConfig struct {
	Port    string
	AppName string

	Logger zap.Logger

	Services *Services
}

func NewServer(c *ServerConfig) *Server {
	return &Server{
		srv: fiber.New(fiber.Config{
			AppName: c.AppName,
		}),
		logger:      c.Logger,
		conf:        c,
		UserService: c.Services.UserService,
		AuthService: c.Services.AuthService,
	}
}

func (s *Server) Run() error {
	err := s.srv.Listen(s.conf.Port)
	if err != nil {
		s.logger.Error("startup server", zap.Error(err))
		return err
	}
	return nil
}

func (s *Server) Build() error {
	s.srv.Use(LoggerMiddleware(s.logger))

	// create superuser
	_, err := s.UserService.GetByLogin("admin")
	if err != nil {
		_, err = s.UserService.Create(serv_dto.User{
			Login:         "admin", // default login
			Password:      "admin", // default password
			IsAdmin:       true,    // admin
			CheckPassword: false,   // no need to check password
		})

		if err != nil {
			s.logger.Error("create superuser", zap.Error(err))
			return err
		}

		s.logger.Info("superuser created")
	}

	s.srv.Get("/swagger/*", SwaggerMiddleware())
	s.srv.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/swagger/")
	})

	route := s.srv.Group("/api/v1/")

	route.Post("/auth/login", LoginHandler(s.AuthService))
	route.Post("/auth/register", RegisterHandler(s.UserService))

	registred := route.Use(CheckAuthorizeMiddleware(s.AuthService))
	registred.Post("/auth/logout", LogoutHandler(s.AuthService))
	registred.Get("/users/me", GetUserSelfHandler(s.UserService))
	registred.Post("/users/me/change-pass", UserSelfChangePasswordHandler(s.UserService))

	return nil
}
