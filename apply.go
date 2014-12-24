package main

import (
	"github.com/codegangsta/cli"
	"github.com/elwinar/rambler/migration"
)

func Apply(c *cli.Context) {
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
		Info.Println("migration table not found")
		err := s.CreateMigrationTable()
		if err != nil {
			Error.Fatalln("unable to create the migration table:", err)
		}
		Info.Println("created")
	}

	applied, err := s.ListAppliedMigrations()
	if err != nil {
		Error.Fatalln("failed to list applied migrations:", err)
	}

	available, err := s.ListAvailableMigrations()
	if err != nil {
		Error.Fatalln("failed to list available migrations:", err)
	}
	
	var i, j int = 0, 0
	for i < len(available) && j < len(applied) {
		if available[i] < applied[j] {
			Error.Fatalln("out of order migration", available[i])
		}

		if available[i] > applied[j] {
			Error.Fatalln("missing mgiration", applied[j])
		}

		i++
		j++
	}

	if j != len(applied) {
		Error.Fatalln("missing mgiration", applied[j])
	}

	for _, v := range available[i:] {
		m, err := migration.NewMigration(Env.Directory, v)
		if err != nil {
			Error.Fatalln("failed to retrieve migration", v, ":", err)
		}
		
		Info.Println("applying", m.Name)

		tx, err := s.StartTransaction()
		if err != nil {
			Error.Fatalln("failed to start transaction:", err)
		}

		for _, statement := range m.Scan("up") {
			Debug.Println(statement)
			_, sqlerr := tx.Exec(statement)
			
			if sqlerr != nil {
				Error.Println("migration failed:", sqlerr)
				txerr := tx.Rollback()
				
				if txerr != nil {
					Error.Fatalln("rollback failed:", txerr)
				}
				
				Error.Fatalln("successfully rolled back")
			}
		}

		err = s.SetMigrationApplied(m.Version, m.Description)
		if err != nil {
			Error.Fatalln("unable to set migration as applied:", err)
		}

		txerr := tx.Commit()
		if txerr != nil {
			Error.Fatalln("commit failed:", txerr)
		}

		if !c.Bool("all") {
			break
		}
	}

	return
}
