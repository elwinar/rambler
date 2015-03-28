package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

func Apply(ctx *cli.Context) {
	err := apply(service, ctx.Bool("all"))
	if err != nil {
		log.Println(err)
	}
}

func apply(service Servicer, all bool) error {
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		err := service.Initialize()
		if err != nil {
			return fmt.Errorf("unable to initialize database: %s", err)
		}
	}

	available, err := service.Available()
	if err != nil {
		return fmt.Errorf("unable to retrieve available migrations: %s", err)
	}

	applied, err := service.Applied()
	if err != nil {
		return fmt.Errorf("unable to retrieve applied migrations: %s", err)
	}

	var i, j = 0, 0
	for i < len(available) && j < len(applied) {
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
		return nil
	}

	for _, migration := range available[i:] {
		err := service.Apply(migration)
		if err != nil {
			return err
		}
		
		if !all {
			return nil
		}
	}

	return nil
}
