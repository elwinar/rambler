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

type StandardService struct {
	driver Driver
	env    configuration.Environment
}

// NewService initialize a new service with the given informations
func NewService(env configuration.Environment) (Service, error) {
	return newService(env, os.Stat, GetDriver)
}

type stater func(string) (os.FileInfo, error)
type driverConstructor func(configuration.Environment) (Driver, error)

func newService(env configuration.Environment, stat stater, newDriver driverConstructor) (*StandardService, error) {
	if _, err := stat(env.Directory); err != nil {
		return nil, ErrUnknownDirectory
	}

	driver, err := newDriver(env)
	if err != nil {
		return nil, ErrUnknownDriver
	}

	return &StandardService{
		driver: driver,
		env:    env,
	}, nil
}

func (s StandardService) MigrationTableExists() (bool, error) {
	return s.driver.MigrationTableExists()
}

func (s StandardService) CreateMigrationTable() error {
	return s.driver.CreateMigrationTable()
}

func (s StandardService) ListAppliedMigrations() ([]uint64, error) {
	return s.driver.ListAppliedMigrations()
}

func (s StandardService) ListAvailableMigrations() ([]uint64, error) {
	return s.driver.ListAvailableMigrations()
}
