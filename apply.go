package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

// Apply available migrations based on the provided context.
func Apply(ctx *cli.Context) {
	err := apply(service, ctx.Bool("all"))
	if err != nil {
		log.Println(err)
	}
}

func apply(service Servicer, all bool) error {
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

	if i == len(available) {
		logger.Info("no migration to apply")
		return nil
	}

	logger.Info("%d migrations to apply", len(available[i:]))
	for _, migration := range available[i:] {
		logger.Info("applying %s", migration.Name)
		err := service.Apply(migration)
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
