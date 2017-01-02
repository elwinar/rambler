package main

import (
	"fmt"
	"os"

	"github.com/bradfitz/slice"
	"github.com/elwinar/rambler/log"
	"github.com/urfave/cli"
)

// Reverse available migrations based on the provided context.
func Reverse(ctx *cli.Context) {
	err := reverse(service, ctx.Bool("all"), logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func reverse(service Servicer, all bool, logger *log.Logger) error {
	logger.Debug("checking database state")
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		return fmt.Errorf("uninitialized database")
	}

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

	logger.Debug("reversing order of applied migrations")
	slice.Sort(applied, func(i, j int) bool {
		return applied[i].Name > applied[j].Name
	})

	logger.Info("%d migrations to reverse", len(applied))
	for _, migration := range applied {
		logger.Info("reversing %s", migration.Name)
		err := service.Reverse(migration)
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
