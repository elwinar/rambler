package migration

import (
	"github.com/elwinar/rambler/migration/driver"
)

type MockConn struct {
	migrationTableExists  func() (bool, error)
	createMigrationTable  func() error
	listAppliedMigrations func() ([]uint64, error)
	setMigrationApplied   func(uint64, string) error
	unsetMigrationApplied   func(uint64) error
	startTransaction      func() (driver.Tx, error)
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

func (c *MockConn) UnsetMigrationApplied(version uint64) error {
	return c.unsetMigrationApplied(version)
}

func (c *MockConn) StartTransaction() (driver.Tx, error) {
	return c.startTransaction()
}
