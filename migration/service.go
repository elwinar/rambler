package migration

import (
	"github.com/elwinar/rambler/configuration"
	"os"
)

// Service gather operations to manipulate migrations table and migrations on
// the filesystem.
type Service interface {
	Driver
}

type service struct {
	driver    Driver
	directory string
}

// NewService initialize a new service with the given informations
func NewService(env configuration.Environment, directory string) (Service, error) {
	return newService(env, directory, os.Stat, GetDriver)
}

type stater func(string) (os.FileInfo, error)
type driverConstructor func(configuration.Environment) (Driver, error)

func newService(env configuration.Environment, directory string, stat stater, newDriver driverConstructor) (*service, error) {
	if _, err := stat(directory); err != nil {
		return nil, ErrUnknownDirectory
	}

	driver, err := newDriver(env)
	if err != nil {
		return nil, ErrUnknownDriver
	}

	return &service{
		driver:    driver,
		directory: directory,
	}, nil
}

func (s service) MigrationTableExists() (bool, error) {
	return s.driver.MigrationTableExists()
}

func (s service) CreateMigrationTable() error {
	return s.driver.CreateMigrationTable()
}
