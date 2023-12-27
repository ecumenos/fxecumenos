package migrations

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"

	// this impost is needed for running migration down
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// this impost is needed for running migration down
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewMigrateUpFunc() func(sourceURL, dbURL string, log *zap.Logger, shutdowner fx.Shutdowner) error {
	return func(sourceURL, dbURL string, log *zap.Logger, shutdowner fx.Shutdowner) error {
		log.Info("command starting...")
		defer func(l *zap.Logger) {
			l.Info("command finished")
		}(log)

		if err := migrateUp(sourceURL, dbURL); err != nil {
			return err
		}

		return shutdowner.Shutdown()

	}
}

func migrateUp(sourceURL, dbURL string) error {
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func NewMigrateDownFunc() func(sourceURL, dbURL string, log *zap.Logger, shutdowner fx.Shutdowner) error {
	return func(sourceURL, dbURL string, log *zap.Logger, shutdowner fx.Shutdowner) error {
		log.Info("command starting...")
		defer func(l *zap.Logger) {
			l.Info("command finished")
		}(log)

		if err := migrateDown(sourceURL, dbURL); err != nil {
			return err
		}

		return shutdowner.Shutdown()

	}
}

func migrateDown(sourceURL, dbURL string) error {
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return err
	}
	if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
