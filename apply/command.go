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
	_ = service
	return err
}
