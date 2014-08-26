package main

import (
	"errors"
	"fmt"
	"github.com/gonuts/commander"
)

const (
	errOutOfOrder        = `migration file %s is out of order`
	queryInsertMigration = `INSERT INTO %s (version, date, description, file) VALUES (%d, NOW(), "%s", "%s")`
	msgDatabaseUpToDate  = `database up-to-date`
)

// migrate run the migration sequence
func migrate(command *commander.Command, args []string) error {
	// Parse the configuration from the command flags
	err := GetConfig(command.Flag)
	if err != nil {
		return err
	}

	// Get the database
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	// Get the migrations rows
	migrations, err := GetMigrations(db)
	if err != nil {
		return err
	}

	// Get the migrations files
	files, err := GetMigrationFiles("./*.sql")
	if err != nil {
		return err
	}

	// Do a parallel iteration over the migrations and files to find the last applied migration,
	// and warn about out-of-order new migrations
	i, j := 0, 0
	for i < len(files) && j < len(migrations) {
		if files[i].Version == migrations[j].Version {
			i++
			j++
			continue
		}

		if files[i].Version < migrations[j].Version {
			return errors.New(fmt.Sprintf(errOutOfOrder, files[i].File))
		}

		if files[i].Version > migrations[j].Version {
			return errors.New(fmt.Sprintf(errMissingMigration, migrations[j].File))
		}
	}

	// If both counters are at the end of the slice, everything is okay
	if i == len(files) && j == len(migrations) {
		fmt.Println(msgDatabaseUpToDate)
		return nil
	}

	// If there more migrations in the database than files, there is files missing
	if j < len(migrations) {
		return errors.New(fmt.Sprintf(errMissingMigration, migrations[j].File))
	}

	// If there is more files than migrations, there is work to do!
	for i < len(files) {
		statements, err := files[i].Up()
		if err != nil {
			return err
		}

		fmt.Println(files[i].File)

		tx, err := db.Beginx()
		if err != nil {
			return err
		}

		for _, statement := range statements {
			_, err := db.Exec(statement)
			if err != nil {
				fmt.Println(err)
				err := tx.Rollback()
				if err != nil {
					return err
				}
				return nil
			}
		}

		err = tx.Commit()
		if err != nil {
			fmt.Println(err)
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return nil
		}

		_, err = db.Exec(fmt.Sprintf(queryInsertMigration, tableMigrations, files[i].Version, files[i].Description, files[i].File))
		if err != nil {
			return err
		}
		
		i++
	}

	fmt.Println(msgDatabaseUpToDate)
	return nil
}
