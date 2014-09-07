package main

import (
	"errors"
	"fmt"
	"github.com/gonuts/commander"
)

const (
	queryDeleteMigration = `DELETE FROM %s WHERE version = %d`
	msgDatabaseRolledBack = `database rolled-back`
)

// rollback reverse the migration sequence
func rollback(command *commander.Command, args []string) error {
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
	i, j := len(files)-1, len(migrations)-1
	
	for i >= 0 && j >= 0 {
		
		// If both migrations match, rollback!
		if files[i].Version == migrations[j].Version {
			
			statements, err := files[i].Down()
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
			
			_, err = db.Exec(fmt.Sprintf(queryDeleteMigration, tableMigrations, files[i].Version))
			if err != nil {
				return err
			}
			
			i--
			j--
			continue
		}
		
		// If file version is higher than migration version, it is out-of-order
		if files[i].Version > migrations[j].Version {
			return errors.New(fmt.Sprintf(errOutOfOrder, files[i].File))
		}
		
		// If file version is lower than migration version, there is a missing file
		if files[i].Version < migrations[j].Version {
			return errors.New(fmt.Sprintf(errMissingMigration, migrations[j].File))
		}
	}
	
	// If there is still files, they are out of order files
	if i >= 0 {
		return errors.New(fmt.Sprintf(errOutOfOrder, files[i].File))
	}
	
	// If there is still migrations, there is files missing
	if j >= 0 {
		return errors.New(fmt.Sprintf(errMissingMigration, migrations[j].File))
	}
	
	fmt.Println(msgDatabaseRolledBack)
	return nil
}
