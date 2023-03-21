package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
	"github.com/petara94/go-auth/cmd/migrations"
	"github.com/petara94/go-auth/internal/config"
	"github.com/petara94/go-auth/internal/transport/http/api"
	"github.com/petara94/go-auth/logger"
	"go.uber.org/zap"
	"log"
)

const (
	appName = "go-auth"
)

var Version = ""

func main() {
	lg, err := logger.Logger()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ReadConfigFromFile("config.yml")
	if err != nil {
		lg.Error("load config failed", zap.Error(err))
	}

	if err = run(cfg, lg); err != nil {
		lg.Error("", zap.Error(err))
	}
}

func run(cfg *config.AppConfig, logger *zap.Logger) error {
	if cfg == nil {
		return ErrNilConfig
	}

	ctx := context.Background()

	poolCfg, err := pgxpool.ParseConfig(cfg.DBConf.URL)
	if err != nil {
		return err
	}

	poolCfg.MaxConns = cfg.DBConf.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		logger.Error("db connection", zap.Error(err))
		return err
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("db ping", zap.Error(err))
		return err
	}

	err = migrations.MigrateDatabase(&cfg.DBConf, *logger)
	if err != nil {
		logger.Error("migrations filed", zap.Error(err))
		return err
	}

	server := api.NewServer(&api.ServerConfig{
		Port:     cfg.API.Port,
		AppName:  appName,
		Services: api.NewServices(ctx, pool, *logger),
		Logger:   *logger,
	})

	err = server.Build()
	if err != nil {
		logger.Error("building server", zap.Error(err))
		return err
	}

	logger.Info("app started", zap.String("version", Version))

	return server.Run()
}
