package driver

import (
	"github.com/elwinar/rambler/configuration"
)

type MockDriver struct {
	new func(configuration.Environment) (Conn, error)
}

func (d *MockDriver) New(env configuration.Environment) (Conn, error) {
	return d.new(env)
}

type MockConn struct {
	migrationTableExists  func() (bool, error)
	createMigrationTable  func() error
	listAppliedMigrations func() ([]uint64, error)
	setMigrationApplied   func(uint64, string) error
	startTransaction      func() (Tx, error)
}

func (c *MockConn) MigrationTableExists() (bool, error) {
	return c.migrationTableExists()
}

func (c *MockConn) CreateMigrationTable() error {
	return c.createMigrationTable()
}

func (c *MockConn) ListAppliedMigrations() ([]uint64, error) {
	return c.listAppliedMigrations()
}

func (c *MockConn) SetMigrationApplied(version uint64, description string) error {
	return c.setMigrationApplied(version, description)
}

func (c *MockConn) StartTransaction() (Tx, error) {
	return c.startTransaction()
}
