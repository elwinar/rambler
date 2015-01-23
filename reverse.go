package main

import (
	"github.com/codegangsta/cli"
	"github.com/elwinar/rambler/migration"
)

func Reverse(c *cli.Context) {
	Env, Debug, Info, Error, err := bootstrap(c)
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

	if len(applied) == 0 {
		return
	}

	var i, j int = len(available) - 1, len(applied) - 1
	for i >= 0 && available[i] > applied[j] {
		i--
	}

	for i >= 0 && j >= 0 {
		if available[i] < applied[j] {
			Error.Fatalln("missing migration", applied[j])
		}

		if available[i] > applied[j] {
			Error.Fatalln("out of order migration", available[i])
		}

		i--
		j--
	}

	if j >= 0 {
		Error.Fatalln("missing migration", applied[j])
	}

	if i >= 0 {
		Error.Fatalln("out of order migration", available[i])
	}

	for i := len(applied) - 1; i >= 0; i-- {
		v := applied[i]

		m, err := migration.NewMigration(Env.Directory, v)
		if err != nil {
			Error.Fatalln("failed to retrieve migration", v, ":", err)
		}

		Info.Println("applying", m.Name)

		statements := m.Scan("down")
		for i := len(statements) - 1; i >= 0; i-- {
			statement := statements[i]
			Debug.Println(statement)
			err := s.Exec(statement)
			if err != nil {
				Error.Fatalln("migration failed:", err)
			}
		}

		err = s.UnsetMigrationApplied(m.Version)
		if err != nil {
			Error.Fatalln("unable to unset migration as applied:", err)
		}

		if !c.Bool("all") {
			break
		}
	}

	return
}
