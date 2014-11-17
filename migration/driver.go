package migration

import (
	"errors"
)

type Driver interface{}

var (
	drivers map[string]Driver
)

func init() {
	drivers = make(map[string]Driver)
}

var (
	ErrDriverAlreadyRegistered = errors.New("driver already registered")
	ErrDriverNotRegistered = errors.New("driver not registered")
)

func RegisterDriver(name string, driver Driver) error {
	return registerDriver(name, driver, drivers)
}

func registerDriver(name string, driver Driver, drivers map[string]Driver) error {
	if _, found := drivers[name]; found {
		return ErrDriverAlreadyRegistered
	}
	
	drivers[name] = driver
	return nil
}

func GetDriver(name string) (Driver, error) {
	return getDriver(name, drivers)
}

func getDriver(name string, drivers map[string]Driver) (Driver, error) {
	driver, found := drivers[name]
	if !found {
		return nil, ErrDriverNotRegistered
	}
	
	return driver, nil
}
