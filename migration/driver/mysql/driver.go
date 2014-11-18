package mysql

import (
	"github.com/elwinar/rambler/migration/driver"
)

func init() {
	driver.Register("mysql", Constructor)
}

func Constructor(options string) (driver.Driver, error) {
	return Driver{}, nil
}

type Driver struct{}

func (d Driver) MigrationTableExists() (bool, error) {
	return false, nil
}

func (d Driver) CreateMigrationTable() (error) {
	return nil
}
