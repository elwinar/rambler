package driver

import (
	"github.com/elwinar/rambler/configuration"
)

type MockDriver struct {
	initialize func(configuration.Environment) error
	migrationTableExists func() (bool, error)
	createMigrationTable func() error
	listAppliedMigrations func() ([]uint64, error)
	startTransaction func() (Tx, error)
}

func (d *MockDriver) Initialize(env configuration.Environment) error {
	return d.initialize(env)
}

func (d *MockDriver) MigrationTableExists() (bool, error) {
	return d.migrationTableExists()
}

func (d *MockDriver) CreateMigrationTable() error {
	return d.createMigrationTable()
}

func (d *MockDriver) ListAppliedMigrations() ([]uint64, error) {
	return d.listAppliedMigrations()
}

func (d *MockDriver) StartTransaction() (Tx, error) {
	return d.startTransaction()
}
