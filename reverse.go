package main

import (
	"github.com/codegangsta/cli"
	"log"
)

func Reverse(ctx *cli.Context) {
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
		log.Fatalln(`no migration table found, nothing to do`)
	}

	applied, err := s.ListAppliedMigrations()
	if err != nil {
		log.Fatalln(`failed to list applied migrations:`, err)
	}

	available := s.ListAvailableMigrations()

	if len(applied) == 0 {
		return
	}

	var i, j int = len(available) - 1, len(applied) - 1
	for i >= 0 && available[i] > applied[j] {
		i--
	}

	for i >= 0 && j >= 0 {
		if available[i] < applied[j] {
			log.Fatalln(`missing migration`, applied[j])
		}

		if available[i] > applied[j] {
			log.Fatalln(`out of order migration`, available[i])
		}

		i--
		j--
	}

	if j >= 0 {
		log.Fatalln(`missing migration`, applied[j])
	}

	if i >= 0 {
		log.Fatalln(`out of order migration`, available[i])
	}

	for i := len(applied) - 1; i >= 0; i-- {
		v := applied[i]

		m, err := NewMigration(env.Directory, v)
		if err != nil {
			log.Fatalln(`failed to retrieve migration`, v, `:`, err)
		}

		log.Println(`applying`, m.Name)

		statements := m.Scan(`down`)
		for i := len(statements) - 1; i >= 0; i-- {
			statement := statements[i]
			log.Println(statement)
			err := s.Exec(statement)
			if err != nil {
				log.Fatalln(`migration failed:`, err)
			}
		}

		err = s.UnsetMigrationApplied(m.Version)
		if err != nil {
			log.Fatalln(`unable to unset migration as applied:`, err)
		}

		if !ctx.Bool(`all`) {
			break
		}
	}

	return
}
