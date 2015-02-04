package driver

import (
	"fmt"
)

var (
	errInvalid           = "not a valid driver"
	errAlreadyRegistered = "driver %s already registered"
	errNotRegistered     = "driver %s not registered"
)

// Driver is the interface used by the program to initialize the database connection.
type Driver interface {
	New(dns, schema string) (Conn, error)
}

var drivers = make(map[string]Driver)

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
func Get(drv, dsn, schema string) (Conn, error) {
	return get(drv, dsn, schema, drivers)
}

func get(drv, dsn, schema string, drivers map[string]Driver) (Conn, error) {
	driver, found := drivers[drv]
	if !found {
		return nil, fmt.Errorf(errNotRegistered, drv)
	}

	conn, err := driver.New(dsn, schema)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
