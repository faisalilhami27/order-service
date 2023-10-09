package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //nolint:revive,nolintlint
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"order-service/config"
)

// Run execute when u need to use auto Migration for database
// Recommended for development and staging
func Run() error {
	log.SetLevel(log.InfoLevel)
	log.Infof("database auto migration: %v", config.Config.Database.AutoMigrate)

	// init database migration
	// use DB master to migrate.
	source := "file://migrations"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Config.Database.Username,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Name,
	)
	m, err := migrate.New(source, dsn)
	if err != nil {
		log.Errorf("error init golang-migrate %s", err)
		return err
	}

	defer func() {
		errSource, errDatabase := m.Close()

		if errSource != nil {
			log.Errorf("error close source golang-migrate %s", errSource)
		}

		if errDatabase != nil {
			log.Errorf("error close database golang-migrate %s", errDatabase)
		}
	}()

	// if migration has no change, ignore error no change
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
