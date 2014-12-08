package apply

import (
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
)

func Command(env configuration.Environment, all bool) error {
	return command(env, all, migration.NewService)
}

type serviceConstructor func(env configuration.Environment) (migration.Service, error)

func command(env configuration.Environment, all bool, newService serviceConstructor) error {
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
	
	_ = applied
	_ = availables
	return err
}
