package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"go.uber.org/zap"
)

type ID uint64

//go:generate mockery --name UserService
type UserService interface {
	Create(u dto.User) (*dto.User, error)
	GetByID(id uint64) (*dto.User, error)
	Get(perPage, page int) ([]*dto.User, error)
	Update(u dto.User) (*dto.User, error)
	Delete(id uint64) error
}

//go:generate mockery --name UserGroupService
type UserGroupService interface {
	Create(u dto.UserGroup) (*dto.UserGroup, error)
	GetByID(id uint64) (*dto.UserGroup, error)
	Get(perPage, page int) ([]*dto.UserGroup, error)
	Update(u dto.UserGroup) (*dto.UserGroup, error)
	Delete(id uint64) error
}

//go:generate mockery --name AuthService
type AuthService interface {
	Login(auth dto.Auth) (*dto.Session, error)
	Get(token string) (*dto.Session, error)
	Logout(session dto.Session) error
}

type Server struct {
	srv  *fiber.App
	conf *ServerConfig

	logger zap.Logger

	UserService      UserService
	UserGroupService UserGroupService
	AuthService      AuthService
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
		logger:           c.Logger,
		conf:             c,
		UserService:      c.Services.UserService,
		UserGroupService: c.Services.UserGroupService,
		AuthService:      c.Services.AuthService,
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

	s.srv.Get("/swagger/*", SwaggerMiddleware())
	s.srv.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/swagger/")
	})

	route := s.srv.Group("/api/v1/")

	userRoute := route.Group("/users/", CheckAuthorizeMiddleware(s.AuthService))
	userRoute.Get("/",
		CheckAdminMiddleware(s.UserService, s.UserGroupService),
		GetUserAllHandler(s.UserService))
	userRoute.Get("/:id", GetUserByIDHandler(s.UserService))
	userRoute.Put("/:id", UpdateUserHandler(s.UserService))
	userRoute.Delete("/:id", DeleteUserHandler(s.UserService))

	userGroupRoute := route.Group("/user-groups/")
	userGroupRoute.Post("/", CreateUserGroupHandler(s.UserGroupService))
	userGroupRoute.Get("/",
		CheckAdminMiddleware(s.UserService, s.UserGroupService),
		GetUserGroupAllHandler(s.UserGroupService))
	userGroupRoute.Get("/:id", GetUserGroupByIDHandler(s.UserGroupService))
	userGroupRoute.Put("/:id", UpdateUserGroupHandler(s.UserGroupService))
	userGroupRoute.Delete("/:id", DeleteUserGroupHandler(s.UserGroupService))

	auth := route.Group("/auth/")
	auth.Post("/register", CreateUserHandler(s.UserService))
	auth.Post("/login", LoginHandler(s.AuthService))
	auth.Post("/logout", LogoutHandler(s.AuthService))

	return nil
}
