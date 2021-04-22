package main

import (
	"fmt"

	"github.com/elwinar/rambler/log"
	"github.com/urfave/cli"
)

// Apply available migrations based on the provided context.
func Apply(ctx *cli.Context) error {
	return apply(service, ctx.Bool("all"), !ctx.Bool("no-save"), ctx.String("migration"), logger)
}

func apply(service Servicer, all, save bool, migration string, logger *log.Logger) error {
	logger.Debug("checking database state")
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		logger.Info("initializing database")
		err := service.Initialize()
		if err != nil {
			return fmt.Errorf("unable to initialize database: %s", err)
		}
	}

	var targets []Migration
	if migration != "" {
		logger.Debug("fetching migration")
		m, err := NewMigration(migration)
		if err != nil {
			return fmt.Errorf("unable to retrieve migration: %s", err)
		}
		targets = append(targets, m)
	} else {
		logger.Debug("fetching available migrations")
		available, err := service.Available()
		if err != nil {
			return fmt.Errorf("unable to retrieve available migrations: %s", err)
		}
		logger.Info("found %d available migrations", len(available))

		logger.Debug("fetching applied migrations")
		applied, err := service.Applied()
		if err != nil {
			return fmt.Errorf("unable to retrieve applied migrations: %s", err)
		}
		logger.Info("found %d applied migrations", len(applied))

		logger.Debug("comparing available and applied migrations")
		var i, j = 0, 0
		for i < len(available) && j < len(applied) {
			logger.Debug("comparing %s and %s", available[i].Name, applied[j].Name)
			if available[i].Name == applied[j].Name {
				i++
				j++
				continue
			}

			if available[i].Name < applied[j].Name {
				return fmt.Errorf("out of order migration: %s", available[i].Name)
			}

			if available[i].Name > applied[j].Name {
				return fmt.Errorf("missing migration: %s", applied[j].Name)
			}
		}

		if j != len(applied) {
			return fmt.Errorf("missing migration: %s", applied[j].Name)
		}

		targets = available[i:]
	}

	if len(targets) == 0 {
		logger.Info("no migration to apply")
		return nil
	}

	logger.Info("%d migrations to apply", len(targets))
	for _, migration := range targets {
		logger.Info("applying %s", migration.Name)
		err := service.Apply(migration, save)
		if err != nil {
			return err
		}

		if !all {
			break
		}
	}

	logger.Info("done")
	return nil
}
