package sqldb

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func MigrateSQLDB(db *gorm.DB, migrationPath, mode string, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	pgdriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), pgdriver)
	if err != nil {
		return err
	}
	return MigrateSchema(m, mode, steps)
}

func MigrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migartionErr error

	switch mode {
	case "up":
		migartionErr = m.Up()
	case "steps":
		if steps == 0 {
			return errors.New("[Database Migration] Step must not be 0")
		}
		migartionErr = m.Steps(steps)
	default:
		return errors.New("[Database Migration] Invalid mode")
	}
	if migartionErr != nil && !errors.Is(migartionErr, migrate.ErrNoChange) {
		return migartionErr
	}
	return nil

}