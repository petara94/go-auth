package api

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petara94/go-auth/internal/repo"
	"github.com/petara94/go-auth/internal/services"
	"go.uber.org/zap"
)

type Services struct {
	UserService      UserService
	UserGroupService UserGroupService
	AuthService      AuthService
}

func NewServices(ctx context.Context, pool *pgxpool.Pool, logger zap.Logger) *Services {
	return &Services{
		UserService:      services.NewUserService(ctx, repo.NewUserStore(ctx, pool), logger),
		UserGroupService: services.NewUserGroupService(ctx, repo.NewUserGroupRepository(ctx, pool), logger),
		AuthService:      services.NewAuthService(ctx, repo.NewSessionStore(ctx, pool), repo.NewUserStore(ctx, pool)),
	}
}
