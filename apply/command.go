package apply

import (
// 	"errors"
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
)

// Command is the `rambler apply` handler
func Command(env configuration.Environment, all bool) error {
	return command(env, all, migration.NewService, migration.NewMigration, Filter, Apply)
}

func command(env configuration.Environment, all bool, newService serviceConstructor, newMigration migrationConstructor, filter filterer, apply applyer) error {
	s, err := newService(env)
	if err != nil {
		return err
	}
	
	exists, err := s.MigrationTableExists()
	if err != nil {
		return err
	}
	
	if !exists {
		err := s.CreateMigrationTable()
		if err != nil {
			return err
		}
	}
	
	applied, err := s.ListAppliedMigrations()
	if err != nil {
		return err
	}
	
	availables, err := s.ListAvailableMigrations()
	if err != nil {
		return err
	}
	
	filtered, err := filter(applied, availables)
	if err != nil {
		return err
	}
	
	for _, v := range filtered {
		m, err := newMigration(env.Directory, v)
		if err != nil {
			return err
		}
		
		tx, err := s.StartTransaction()
		if err != nil {
			return err
		}
		
		sqlerr, txerr := apply(m, tx)
		if sqlerr != nil && txerr != nil { // Rollback error
			return errors.New(sqlerr.Error() + "; " + txerr.Error())
		} else if sqlerr != nil { // Rollback success
			return sqlerr
		} else if txerr != nil { // Commit error
			return txerr
		}
		
		if !all {
			break
		}
	}
	
	return nil
}
