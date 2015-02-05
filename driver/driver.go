package driver

import (
	"fmt"
)

// Driver is the interface used by the program to initialize the database connection.
type Driver interface {
	New(dns, schema string) (Conn, error)
}

var drivers = make(map[string]Driver)

// Register register a driver
func Register(name string, driver Driver) error {
	if _, found := drivers[name]; found {
		return fmt.Errorf(`driver "%s" already registered`, name)
	}

	if driver == nil {
		return fmt.Errorf(`not a valid driver`)
	}

	drivers[name] = driver
	return nil
}

// Get initialize a driver from the given environment
func Get(drv, dsn, schema string) (Conn, error) {
	driver, found := drivers[drv]
	if !found {
		return nil, fmt.Errorf(`driver "%s" not registered`, drv)
	}

	conn, err := driver.New(dsn, schema)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
