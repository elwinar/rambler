package apply

/*
import (
	"errors"
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
)

// Command is the `rambler apply` handler
func Command(env configuration.Environment, all bool) error {
	return command(env, all, migration.NewService, Filter, func(path string, version uint64) (scanner, error) {
		return migration.NewMigration(path, version)
	}, Apply)
}

type serviceConstructor func(configuration.Environment) (migration.Service, error)
type filter func([]uint64, []uint64) ([]uint64, error)
type migrationConstructor func(string, uint64) (scanner, error)
type applier func(scanner, txer) (error, error)

func command(env configuration.Environment, all bool, newService serviceConstructor, f filter, newMigration migrationConstructor, a applier) error {
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
	
	filtered, err := f(applied, availables)
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
		
		sqlerr, txerr := a(m, tx)
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
*/
