package main

import (
	"fmt"

	"github.com/bradfitz/slice"
	"github.com/elwinar/rambler/log"
	"github.com/urfave/cli"
)

// Reverse available migrations based on the provided context.
func Reverse(ctx *cli.Context) error {
	return reverse(service, ctx.Bool("all"), !ctx.Bool("no-save"), ctx.String("migration"), logger)
}

func reverse(service Servicer, all, save bool, migration string, logger *log.Logger) error {
	logger.Debug("checking database state")
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		return fmt.Errorf("uninitialized database")
	}

	var targets []*Migration

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

		if len(applied) == 0 {
			logger.Info("no migration to reverse")
			return nil
		}

		logger.Debug("rewinding to first applied migration")
		var i, j = len(available) - 1, len(applied) - 1
		for i >= 0 && j >= 0 && available[i].Name > applied[j].Name {
			i--
		}

		logger.Debug("comparing available and applied migrations")
		for i >= 0 && j >= 0 {
			logger.Debug("comparing %s and %s", available[i].Name, applied[j].Name)
			if available[i].Name == applied[j].Name {
				i--
				j--
				continue
			}

			if available[i].Name < applied[j].Name {
				return fmt.Errorf("missing migration: %s", applied[j].Name)
			}

			if available[i].Name > applied[j].Name {
				return fmt.Errorf("out of order migration: %s", available[i].Name)
			}
		}

		if i >= 0 {
			return fmt.Errorf("out of order migration: %s", available[i].Name)
		}

		if j >= 0 {
			return fmt.Errorf("missing migration: %s", applied[j].Name)
		}

		targets = applied
	}

	logger.Debug("reversing order of applied migrations")
	slice.Sort(targets, func(i, j int) bool {
		return targets[i].Name > targets[j].Name
	})

	logger.Info("%d migrations to reverse", len(targets))
	for _, migration := range targets {
		logger.Info("reversing %s", migration.Name)
		err := service.Reverse(migration, save)
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
