package driver

import (
	"database/sql"
	"fmt"
	"github.com/elwinar/rambler/configuration"
)

var (
	errInvalid           = "not a valid driver"
	errAlreadyRegistered = "driver %s already registered"
	errNotRegistered     = "driver %s not registered"
)

// Tx is the interface for an SQL transaction as used by rambler.
type Tx interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

// Driver is the interface used by the program to initialize the database connection.
type Driver interface {
	New(configuration.Environment) (Conn, error)
}

var drivers = make(map[string]Driver)

// Conn is the interface used by the program to manipulate the migration table.
type Conn interface {
	MigrationTableExists() (bool, error)
	CreateMigrationTable() error
	ListAppliedMigrations() ([]uint64, error)
	SetMigrationApplied(version uint64, description string) error
	UnsetMigrationApplied(version uint64) error
	StartTransaction() (Tx, error)
}

// Register register a driver
func Register(name string, driver Driver) error {
	return register(name, driver, drivers)
}

func register(name string, driver Driver, drivers map[string]Driver) error {
	if _, found := drivers[name]; found {
		return fmt.Errorf(errAlreadyRegistered, name)
	}

	if driver == nil {
		return fmt.Errorf(errInvalid)
	}

	drivers[name] = driver
	return nil
}

// Get initialize a driver from the given environment
func Get(env configuration.Environment) (Conn, error) {
	return get(env, drivers)
}

func get(env configuration.Environment, drivers map[string]Driver) (Conn, error) {
	driver, found := drivers[env.Driver]
	if !found {
		return nil, fmt.Errorf(errNotRegistered, env.Driver)
	}

	conn, err := driver.New(env)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
