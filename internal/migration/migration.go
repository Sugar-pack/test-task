// Package migrations - database migrations utilities
package migration

import (
	"context"

	migrate "github.com/rubenv/sql-migrate"

	"github.com/Sugar-pack/test-task/internal/config"
	"github.com/Sugar-pack/test-task/internal/logging"
)

// Apply applies database migrations.
func Apply(ctx context.Context, conf *config.DB) error {
	logger := logging.FromContext(ctx)

	dbConn, err := Connect(ctx, conf)
	if err != nil {
		logger.WithError(err).Error("db connect failed")

		return err
	}
	defer func() {
		clErr := dbConn.Close()
		if clErr != nil {
			logger.WithError(clErr).Error("closing db connection is failed")
		}
	}()

	logger.Trace("applying migrations")
	migrationDirection := migrate.Up
	migrationCount := -1
	migrate.SetTable(conf.MigrationTable)
	count, err := migrate.ExecMax(dbConn.DB, "postgres",
		&migrate.FileMigrationSource{Dir: conf.MigrationDirPath},
		migrationDirection, migrationCount,
	)
	if err != nil {
		logger.WithError(err).Error("apply migration failed")

		return err
	}

	logger.WithField("count", count).Info("migrations applied")

	if err = Disconnect(ctx, dbConn); err != nil {
		logger.WithError(err).Error("disconnect failed")

		return err
	}

	return err
}
