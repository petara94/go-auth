package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/petara94/go-auth/internal/config"
	"go.uber.org/zap"
)

func MigrateDatabase(cfg *config.DBConf, logger zap.Logger) error {
	m, err := migrate.New("file://migrations", cfg.URL)
	if err != nil {
		logger.Error("migrations build filed", zap.Error(err))
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error("migrations filed", zap.Error(err))
		return err
	}

	if err != migrate.ErrNoChange {
		logger.Info("applied migrations")
	} else {
		logger.Info("migrations aren't needed")
	}

	return nil
}
