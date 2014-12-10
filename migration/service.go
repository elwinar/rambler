package migration

import (
	"github.com/elwinar/rambler/configuration"
	"os"
)

// Service gather operations to manipulate migrations table and migrations on
// the filesystem.
type Service interface {
	Driver
	ListAvailableMigrations() ([]uint64, error)
}

type service struct {
	Driver
	env    configuration.Environment
}

// NewService initialize a new service with the given informations
func NewService(env configuration.Environment) (Service, error) {
	return newService(env, os.Stat, GetDriver)
}

type stater func(string) (os.FileInfo, error)
type driverConstructor func(configuration.Environment) (Driver, error)

func newService(env configuration.Environment, stat stater, newDriver driverConstructor) (*service, error) {
	if _, err := stat(env.Directory); err != nil {
		return nil, ErrUnknownDirectory
	}

	driver, err := newDriver(env)
	if err != nil {
		return nil, ErrUnknownDriver
	}

	return &service{
		Driver: driver,
		env:    env,
	}, nil
}

func (s service) ListAvailableMigrations() ([]uint64, error) {
	return listAvailableMigrations(s.env)
}

func listAvailableMigrations(env configuration.Environment) ([]uint64, error) {
	
	return nil, nil
}
