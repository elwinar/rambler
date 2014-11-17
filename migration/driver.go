package migration

import (
	"errors"
)

var (
	drivers map[string]Driver
)

func init() {
	drivers = make(map[string]Driver)
}

var (
	ErrDriverAlreadyRegistered = errors.New("driver is already registered")
)

type Driver interface{}

func Register(name string, driver Driver) error {
	return register(name, driver, drivers)
}

func register(name string, driver Driver, drivers map[string]Driver) error {
	if _, found := drivers[name]; found {
		return ErrDriverAlreadyRegistered
	}
	
	drivers[name] = driver
	return nil
}
