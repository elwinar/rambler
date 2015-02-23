package main

import (
	"github.com/codegangsta/cli"
	"log"
)

func Apply(ctx *cli.Context) {
	cfg, err := Load(ctx.GlobalString(`configuration`))
	if err != nil {
		log.Fatalln(`unable to load configuration file:`, err)
	}

	env, err := cfg.Env(ctx.GlobalString(`environment`))
	if err != nil {
		log.Fatalln(`unable to load requested environment:`, err)
	}

	s, err := NewService(env)
	if err != nil {
		log.Fatalln(`unable to initialize the migration service:`, err)
	}

	exists, err := s.MigrationTableExists()
	if err != nil {
		log.Fatalln(`failed to look for migration table:`, err)
	}

	if !exists {
		log.Println(`migration table not found`)
		err := s.CreateMigrationTable()
		if err != nil {
			log.Fatalln(`unable to create the migration table:`, err)
		}
		log.Println(`created`)
	}

	applied, err := s.ListAppliedMigrations()
	if err != nil {
		log.Fatalln(`failed to list applied migrations:`, err)
	}

	available, err := s.ListAvailableMigrations()
	if err != nil {
		log.Fatalln(`failed to list available migrations:`, err)
	}

	var i, j int = 0, 0
	for i < len(available) && j < len(applied) {
		if available[i] < applied[j] {
			log.Fatalln(`out of order migration`, available[i])
		}

		if available[i] > applied[j] {
			log.Fatalln(`missing mgiration`, applied[j])
		}

		i++
		j++
	}

	if j != len(applied) {
		log.Fatalln(`missing mgiration`, applied[j])
	}

	for _, v := range available[i:] {
		m, err := NewMigration(env.Directory, v)
		if err != nil {
			log.Fatalln(`failed to retrieve migration`, v, `:`, err)
		}

		log.Println(`applying`, m.Name)

		statements := m.Scan(`up`)

		for _, statement := range statements {
			log.Println(statement)
			err := s.Exec(statement)
			if err != nil {
				log.Fatalln(`migration failed:`, err)
			}
		}

		err = s.SetMigrationApplied(m.Version, m.Description)
		if err != nil {
			log.Fatalln(`unable to set migration as applied:`, err)
		}

		if !ctx.Bool(`all`) {
			break
		}
	}

	return
}
