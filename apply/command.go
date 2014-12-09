package apply

import (
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
)

func Command(env configuration.Environment, all bool) error {
	return command(env, all, migration.NewService, Filter)
}

type serviceConstructor func(env configuration.Environment) (migration.Service, error)
type filter func([]uint64, []uint64) ([]uint64, error)

func command(env configuration.Environment, all bool, newService serviceConstructor, f filter) error {
	service, err := newService(env)
	if err != nil {
		return err
	}
	
	exists, err := service.MigrationTableExists()
	if err != nil {
		return err
	}
	
	if !exists {
		err := service.CreateMigrationTable()
		if err != nil {
			return err
		}
	}
	
	applied, err := service.ListAppliedMigrations()
	if err != nil {
		return err
	}
	
	availables, err := service.ListAvailableMigrations()
	if err != nil {
		return err
	}
	
	filtered, err := f(applied, availables)
	if err != nil {
		return err
	}
	
	_ = filtered
	
	return err
}
