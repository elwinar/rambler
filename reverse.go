package main

import (
	"github.com/codegangsta/cli"
	"github.com/elwinar/rambler/migration"
)

func Reverse(c *cli.Context) {
	Env, _, Info, Error, err := bootstrap(c)
	if err != nil {
		Error.Fatalln("unable to load configuration file:", err)
	}

	s, err := migration.NewService(*Env)
	if err != nil {
		Error.Fatalln("unable to initialize the migration service:", err)
	}

	exists, err := s.MigrationTableExists()
	if err != nil {
		Error.Fatalln("failed to look for migration table:", err)
	}
	
	if !exists {
		Error.Fatalln("no migration table found, nothing to do")
	}

	applied, err := s.ListAppliedMigrations()
	if err != nil {
		Error.Fatalln("failed to list applied migrations:", err)
	}

	available, err := s.ListAvailableMigrations()
	if err != nil {
		Error.Fatalln("failed to list available migrations:", err)
	}

	_ = Info
	_ = applied
	_ = available

	return
}
