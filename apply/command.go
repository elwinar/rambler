package apply

import (
	"fmt"
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
		return fmt.Errorf(errInitializeNewService, err.Error())
	}

	exists, err := s.MigrationTableExists()
	if err != nil {
		return fmt.Errorf(errMigrationTableCheck, err.Error())
	}

	if !exists {
		err := s.CreateMigrationTable()
		if err != nil {
			return fmt.Errorf(errCreateMigrationTable, err.Error())
		}
	}

	applied, err := s.ListAppliedMigrations()
	if err != nil {
		return fmt.Errorf(errListApplied, err.Error())
	}

	availables, err := s.ListAvailableMigrations()
	if err != nil {
		return fmt.Errorf(errListAvailable, err.Error())
	}

	filtered, err := filter(applied, availables)
	if err != nil {
		return fmt.Errorf(errFilter, err.Error())
	}

	for _, v := range filtered {
		m, err := newMigration(env.Directory, v)
		if err != nil {
			return fmt.Errorf(errNewMigration, v, err.Error())
		}

		tx, err := s.StartTransaction()
		if err != nil {
			return fmt.Errorf(errStartTransactionFailed, err.Error())
		}

		sqlerr, txerr := apply(m.Scan("up"), tx)
		if sqlerr != nil && txerr != nil {
			return fmt.Errorf(errRollbackFailed, sqlerr.Error(), txerr.Error())
		} else if sqlerr != nil {
			return fmt.Errorf(errMigrationError, sqlerr.Error())
		} else if txerr != nil {
			return fmt.Errorf(errCommitFailed, txerr.Error())
		}

		err = s.SetMigrationApplied(m.Version, m.Description)
		if err != nil {
			return fmt.Errorf(errSetMigrationApplied, m.Version, err.Error())
		}

		if !all {
			break
		}
	}

	return nil
}
