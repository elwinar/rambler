package main

import (
	"fmt"
	"github.com/bradfitz/slice"
	"github.com/codegangsta/cli"
	"log"
)

func Reverse(ctx *cli.Context) {
	err := reverse(service, ctx.Bool("all"))
	if err != nil {
		log.Println(err)
	}
}

func reverse(service Servicer, all bool) error {
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		return fmt.Errorf("uninitialized database")
	}

	available, err := service.Available()
	if err != nil {
		return fmt.Errorf("unable to retrieve available migrations: %s", err)
	}

	applied, err := service.Applied()
	if err != nil {
		return fmt.Errorf("unable to retrieve applied migrations: %s", err)
	}
	
	if len(applied) == 0 {
		return nil
	}

	var i, j = len(available) - 1, len(applied) - 1
	for i >= 0 && j >= 0 && available[i].Name > applied[j].Name {
		i--
	}
	
	for i >= 0 && j >= 0 {
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
	
	slice.Sort(applied, func(i, j int) bool {
		return applied[i].Name > applied[j].Name
	})

	for _, migration := range applied {
		err := service.Reverse(migration)
		if err != nil {
			return err
		}
		
		if !all {
			return nil
		}
	}

	return nil
}
